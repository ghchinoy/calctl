package cmd

import (
	"fmt"

	"github.com/ghchinoy/calctl/internal/calendar"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [event-id]",
	Short: "Get the details of a specific event.",
	Long:  `Get the details of a specific event by providing its ID.`, 
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eventID := args[0]
		event, err := calendar.GetEvent(calendarSvc, eventID)
		if err != nil {
			return err
		}

		fmt.Printf("Summary: %s\n", event.Summary)
		fmt.Printf("Start: %s\n", event.Start.DateTime)
		fmt.Printf("End: %s\n", event.End.DateTime)
		fmt.Printf("Description: %s\n", event.Description)
		fmt.Printf("Hangout Link: %s\n", event.HangoutLink)
		fmt.Printf("ID: %s\n", event.Id)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
