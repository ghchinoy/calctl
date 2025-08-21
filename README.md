# calctl

A command-line tool and MCP server for interacting with Google Calendar.

## Installation

To build the tool from source, you need to have Go installed. Then run the following commands:

```bash
git clone https://github.com/ghchinoy/calctl
cd calctl
go build ./...
```

## Configuration

`calctl` requires Google Cloud credentials to access your calendar data.

1.  **Create a Google Cloud Project:** If you don't have one already, create a project in the [Google Cloud Console](https://console.cloud.google.com/).
2.  **Enable the Google Calendar API:**
    ```bash
    gcloud services enable calendar-json.googleapis.com
    ```
3.  **Create OAuth 2.0 Credentials:**
    *   Go to the [Credentials page](https://console.cloud.google.com/apis/credentials) in the Cloud Console.
    *   Click "Create Credentials" and choose "OAuth client ID".
    *   Select "Desktop app" for the application type.
    *   After creation, download the JSON file. It will be named something like `client_secret_xxxxxxxx.json`.
4.  **Set up the credentials for `calctl`:**
    You can either use the `--secret-file` flag with each command:
    ```bash
    ./calctl --secret-file /path/to/your/client_secret.json week
    ```
    Or, you can set the `CALENDAR_SECRETS` environment variable:
    ```bash
    export CALENDAR_SECRETS=/path/to/your/client_secret.json
    ./calctl week
    ```

## Usage

### CLI Commands

*   **`week [date]`**: Display events for the current week. Can optionally take a date in `YYYY-MM-DD` format.
    ```bash
    ./calctl week
    ./calctl week 2025-08-25
    ```
*   **`week --analyze`**: Provides an analysis of the week's events, showing overlaps and your response status.
    ```bash
    ./calctl week --analyze
    ```
*   **`get [event-id]`**: Get the details of a specific event.
    ```bash
    ./calctl get <event_id_from_week_command>
    ```

### MCP Server

You can also run `calctl` as an MCP server to expose its functionality as tools.

*   **To run over standard I/O:**
    ```bash
    ./calctl --mcp
    ```
    By default, logs will be written to stderr. To write logs to a file, use the `--logfile` flag:
    ```bash
    ./calctl --mcp --logfile calctl.log
    ```
*   **To run over HTTP:**
    ```bash
    ./calctl --mcp-http :8080
    ```

The following tools are available:

*   `get_weekly_calendar(current_date: str)`: Get the user's calendar events for the week of the given date (YYYY-MM-DD). If the date is omitted, it uses the current week.
*   `get_event_details(event_id: str)`: Get the details of a specific event by its ID.
