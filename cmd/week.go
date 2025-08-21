package cmd

import (
	"fmt"
	"time"

	"github.com/ghchinoy/calctl/internal/calendar"
	"github.com/spf13/cobra"
)

var analyze bool

// weekCmd represents the week command
var weekCmd = &cobra.Command{
	Use:   "week [date]",
	Short: "Displays events for the current week, or the week of the given date.",
	Long:  `Displays events for the current week (Sunday to Saturday). You can optionally provide a date in YYYY-MM-DD format to see events for that week.`, 
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var dateStr string
		if len(args) > 0 {
			dateStr = args[0]
		}

	
events, err := calendar.GetWeeklyEvents(calendarSvc, dateStr)
		if err != nil {
			return err
		}

		if analyze {
			report, err := calendar.AnalyzeEvents(events)
			if err != nil {
				return err
			}
			fmt.Println(report.FormatReport())
			return nil
		}

		if len(events.Items) == 0 {
			fmt.Println("No upcoming events found for this week.")
		} else {
			fmt.Println("Events for this week:")
			for _, item := range events.Items {
				date := item.Start.DateTime
				if date == "" {
					date = item.Start.Date
				}
				t, err := time.Parse(time.RFC3339, date)
				if err != nil {
					// try parsing as a date
					t, err = time.Parse("2006-01-02", date)
					if err != nil {
						fmt.Printf("Unable to parse date: %s\n", date)
						continue
					}
				}

				fmt.Printf("- %s (%s) [%s]\n", item.Summary, t.Format("Mon, Jan 2 3:04 PM"), item.Id)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(weekCmd)
	weekCmd.Flags().BoolVar(&analyze, "analyze", false, "analyze the week's events")
}
