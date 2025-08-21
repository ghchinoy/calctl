package calendar

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

// AnalysisReport contains the results of a calendar analysis.

type AnalysisReport struct {
	OverlappingEvents [][]*calendar.Event
	EventsByStatus    map[string][]*calendar.Event
}

// AnalyzeEvents performs an analysis of the given events.
func AnalyzeEvents(events *calendar.Events) (*AnalysisReport, error) {
	report := &AnalysisReport{
		EventsByStatus: make(map[string][]*calendar.Event),
	}

	// Sort events by start time for overlap detection.
	sort.Slice(events.Items, func(i, j int) bool {

		ti, err := time.Parse(time.RFC3339, events.Items[i].Start.DateTime)
		if err != nil {
			return false
		}
		tj, err := time.Parse(time.RFC3339, events.Items[j].Start.DateTime)
		if err != nil {
			return false
		}
		return ti.Before(tj)
	})

	// Detect overlapping events.
	for i := 0; i < len(events.Items)-1; i++ {
		event1 := events.Items[i]
		event2 := events.Items[i+1]

		end1, err := time.Parse(time.RFC3339, event1.End.DateTime)
		if err != nil {
			continue
		}
		start2, err := time.Parse(time.RFC3339, event2.Start.DateTime)
		if err != nil {
			continue
		}

		if end1.After(start2) {
			report.OverlappingEvents = append(report.OverlappingEvents, []*calendar.Event{event1, event2})
		}
	}

	// Categorize events by status.
	for _, event := range events.Items {
		for _, attendee := range event.Attendees {
			if attendee.Self {
				report.EventsByStatus[attendee.ResponseStatus] = append(report.EventsByStatus[attendee.ResponseStatus], event)
				break
			}
		}
	}

	return report, nil
}

// FormatReport formats the analysis report into a string.
func (r *AnalysisReport) FormatReport() string {
	var builder strings.Builder

	builder.WriteString("\n--- Calendar Analysis Report ---\n")

	if len(r.OverlappingEvents) > 0 {
		builder.WriteString("\nOverlapping Meetings:\n")
		for _, pair := range r.OverlappingEvents {
			builder.WriteString(fmt.Sprintf("  - '%s' and '%s'\n", pair[0].Summary, pair[1].Summary))
		}
	} else {
		builder.WriteString("\nNo overlapping meetings found.\n")
	}

	if len(r.EventsByStatus) > 0 {
		builder.WriteString("\nMeetings by Your Status:\n")
		for status, events := range r.EventsByStatus {
			builder.WriteString(fmt.Sprintf("  %s:\n", strings.Title(status)))
			for _, event := range events {
				builder.WriteString(fmt.Sprintf("    - %s\n", event.Summary))
			}
		}
	}

	builder.WriteString("\n--- End of Report ---\n")

	return builder.String()
}
