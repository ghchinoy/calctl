package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ghchinoy/calctl/internal/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	// noBrowserAuth is a flag to disable opening the browser for authentication.
	noBrowserAuth bool
	// client is the HTTP client used for all API calls.
	client *http.Client
	// calendarSvc is the Google Calendar service client.
	calendarSvc *calendar.Service
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "calctl",
	Short: "A CLI for Google Calendar.",
	Long:  `calctl is a command-line tool for interacting with your Google Calendar.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		secretFile := viper.GetString("secret-file")
		if secretFile == "" {
			return fmt.Errorf("client secret file not set. Please use the --secret-file flag or set the CALENDAR_SECRETS environment variable")
		}

		ctx := context.Background()
		var err error
		client, err = auth.NewOAuthClient(ctx, secretFile, noBrowserAuth)
		if err != nil {
			return fmt.Errorf("could not create oauth client: %w", err)
		}

		calendarSvc, err = calendar.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return fmt.Errorf("could not create calendar service: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("secret-file", "", "path to your client secrets file")
	rootCmd.PersistentFlags().BoolVar(&noBrowserAuth, "no-browser-auth", false, "do not open a browser for authentication")
	rootCmd.PersistentFlags().Bool("mcp", false, "enable MCP server mode over stdio")
	rootCmd.PersistentFlags().String("mcp-http", "", "enable MCP server mode over HTTP at the given address")
	rootCmd.PersistentFlags().String("logfile", "", "path to log file for MCP server")
	viper.BindPFlag("secret-file", rootCmd.PersistentFlags().Lookup("secret-file"))
	viper.BindEnv("secret-file", "CALENDAR_SECRETS")
	viper.BindPFlag("mcp", rootCmd.PersistentFlags().Lookup("mcp"))
	viper.BindPFlag("mcp-http", rootCmd.PersistentFlags().Lookup("mcp-http"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()
}
