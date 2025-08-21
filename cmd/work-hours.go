package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ghchinoy/calctl/internal/calendar"
	"github.com/spf13/cobra"
)

// workHoursCmd represents the work-hours command
var workHoursCmd = &cobra.Command{
	Use:   "work-hours",
	Short: "Displays your configured working hours.",
	Long:  `Displays your configured working hours as set in your Google Calendar settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		setting, err := calendar.GetWorkingHours(calendarSvc)
		if err != nil {
			return err
		}

		type workingHoursProperties struct {
			DaysOfWeek []string `json:"daysOfWeek"`
			StartTime  string   `json:"startTime"`
			EndTime    string   `json:"endTime"`
		}

		if setting.Value == "" {
			fmt.Println("No working hours set in your calendar.")
			return nil
		}

		var wh workingHoursProperties
		if err := json.Unmarshal([]byte(setting.Value), &wh); err != nil {
			return fmt.Errorf("could not parse working hours: %w", err)
		}

		fmt.Println("Your configured working hours are:")
		fmt.Printf("  Days: %v\n", wh.DaysOfWeek)
		fmt.Printf("  Start Time: %s\n", wh.StartTime)
		fmt.Printf("  End Time: %s\n", wh.EndTime)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(workHoursCmd)
}
