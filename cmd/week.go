package cmd

import (
	"fmt"
	"time"

	"github.com/ghchinoy/calctl/internal/calendar"
	"github.com/spf13/cobra"
)

// weekCmd represents the week command
var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "Displays events for the current week.",
	Long:  `Displays events for the current week (Sunday to Saturday).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		events, err := calendar.GetWeeklyEvents(calendarSvc)
		if err != nil {
			return err
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
}
