package mcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ghchinoy/calctl/internal/auth"
	"github.com/ghchinoy/calctl/internal/calendar"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	googlecalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// getCalendarSvc creates a new Google Calendar service client.
func getCalendarSvc(ctx context.Context) (*googlecalendar.Service, error) {
	viper.AutomaticEnv()
	secretFile := viper.GetString("secret-file")
	if secretFile == "" {
		return nil, fmt.Errorf("client secret file not set. Please use the --secret-file flag or set the CALENDAR_SECRETS environment variable")
	}
	noBrowserAuth := viper.GetBool("no-browser-auth")
	client, err := auth.NewOAuthClient(ctx, secretFile, noBrowserAuth)
	if err != nil {
		return nil, fmt.Errorf("could not create oauth client: %w", err)
	}
	calendarSvc, err := googlecalendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("could not create calendar service: %w", err)
	}
	return calendarSvc, nil
}

// GetEventArgs defines the arguments for the get_event_details tool.
type GetEventArgs struct {
	EventID string `json:"event_id"`
}

// GetWeeklyCalendarArgs defines the arguments for the get_weekly_calendar tool.
type GetWeeklyCalendarArgs struct {
	CurrentDate string `json:"current_date,omitempty"`
}

// Start starts the MCP server.
func Start(rootCmd *cobra.Command, httpAddr string) error {
	server := mcp.NewServer(&mcp.Implementation{Name: "calctl"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_weekly_calendar",
		Description: "Get the user's calendar events for the current week.",
	}, getWeeklyCalendarHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_event_details",
		Description: "Get the details of a specific event by its ID.",
	}, getEventDetailsHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "analyze_weekly_calendar",
		Description: "Analyzes the user's calendar events for the current week.",
	}, analyzeWeeklyCalendarHandler)

	if httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("MCP handler listening at %s", httpAddr)
		return http.ListenAndServe(httpAddr, handler)
	}

	// MCP over stdio
	logfile := viper.GetString("logfile")
	var logDest io.Writer
	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer f.Close()
		logDest = f
	} else {
		logDest = os.Stderr
	}
	log.SetOutput(logDest)

	t := mcp.NewLoggingTransport(mcp.NewStdioTransport(), logDest)
	if err := server.Run(context.Background(), t); err != nil {
		log.Printf("Server failed: %v", err)
		return err
	}

	return nil
}

// getWeeklyCalendarHandler is the handler for the get_weekly_calendar tool.
func getWeeklyCalendarHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetWeeklyCalendarArgs]) (*mcp.CallToolResultFor[any], error) {
	calendarSvc, err := getCalendarSvc(ctx)
	if err != nil {
		return nil, err
	}

	events, err := calendar.GetWeeklyEvents(calendarSvc, params.Arguments.CurrentDate)
	if err != nil {
		return nil, err
	}

	var output string
	if len(events.Items) == 0 {
		output = "No upcoming events found for this week."
	} else {
		output = "Events for this week:\n"
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			t, err := time.Parse(time.RFC3339, date)
			if err != nil {
				t, err = time.Parse("2006-01-02", date)
				if err != nil {
					output += fmt.Sprintf("- %s (Unable to parse date: %s)\n", item.Summary, date)
					continue
				}
			}
			output += fmt.Sprintf("- %s (%s) [%s]\n", item.Summary, t.Format("Mon, Jan 2 3:04 PM"), item.Id)
		}
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil
}

// getEventDetailsHandler is the handler for the get_event_details tool.
func getEventDetailsHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetEventArgs]) (*mcp.CallToolResultFor[any], error) {
	if params.Arguments.EventID == "" {
		return nil, fmt.Errorf("event_id is a required argument")
	}

	calendarSvc, err := getCalendarSvc(ctx)
	if err != nil {
		return nil, err
	}

	event, err := calendar.GetEvent(calendarSvc, params.Arguments.EventID)
	if err != nil {
		return nil, err
	}

	output := fmt.Sprintf("Summary: %s\nStart: %s\nEnd: %s\nDescription: %s\nHangout Link: %s\nID: %s\n",
		event.Summary, event.Start.DateTime, event.End.DateTime, event.Description, event.HangoutLink, event.Id)

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil
}

// analyzeWeeklyCalendarHandler is the handler for the analyze_weekly_calendar tool.
func analyzeWeeklyCalendarHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetWeeklyCalendarArgs]) (*mcp.CallToolResultFor[any], error) {
	calendarSvc, err := getCalendarSvc(ctx)
	if err != nil {
		return nil, err
	}

	events, err := calendar.GetWeeklyEvents(calendarSvc, params.Arguments.CurrentDate)
	if err != nil {
		return nil, err
	}

	report, err := calendar.AnalyzeEvents(events)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: report.FormatReport()},
		},
	}, nil
}
