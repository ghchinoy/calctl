package cmd

import (
	"fmt"

	"github.com/ghchinoy/calctl/internal/calendar"
	"github.com/spf13/cobra"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Displays all your calendar settings.",
	Long:  `Displays all your calendar settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := calendar.ListSettings(calendarSvc)
		if err != nil {
			return err
		}

		fmt.Println("Your calendar settings:")
		for _, item := range settings.Items {
			fmt.Printf("- %s: %s\n", item.Id, item.Value)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}
