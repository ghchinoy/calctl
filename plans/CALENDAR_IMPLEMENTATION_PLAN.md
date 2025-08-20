# Calctl Implementation Plan

## 1. High-Level Goal

Create a command-line tool (`calctl`) and an associated MCP (Model Context Protocol) server to interact with a user's Google Calendar. The initial features will be to list events for the current week and get details for a specific event.

**Resources**
* /Users/ghchinoy/dev/github/mcp/go-sdk - a local, read-only copy of https://github.com/modelcontextprotocol/go-sdk
* /Users/ghchinoy/dev/github/google-api-go-client/calendar/v3 - a local, read-only copy of the Google Calendar Go SDK

## 2. Core Components & Architecture

The project will be structured similarly to the existing `drivectl` tool, maintaining a clean separation between command-line parsing, core application logic, and API interactions.

*   **`main.go`**: Application entry point.
*   **`cmd/`**: Cobra command definitions.
*   **`internal/`**: Internal packages containing the core logic.
*   **`mcp/`**: MCP server implementation.
*   **`plans/`**: Directory for implementation and test plans.

## 3. Implementation Tasks

### Phase 1: Project Setup & Authentication

*   [x] Create project directories (`cmd`, `internal/auth`, `internal/calendar`, `mcp`).
*   [x] Initialize Go Module (`go.mod`).
*   [x] Add dependencies (`cobra`, `viper`, `google.golang.org/api/calendar/v3`, `golang.org/x/oauth2/google`, `mcp-sdk`).
*   [x] Implement Authentication (`internal/auth/auth.go`) by adapting `drivectl`'s auth flow.
*   [x] Create `main.go` entry point.
*   [x] **Milestone Verification:**
    *   [x] Build the project (`go build ./...`).
    *   [x] Request user to verify by attempting to run the base command, which should trigger the auth flow.

### Phase 2: CLI Commands

*   [x] Create `cmd/root.go` and initialize the calendar service in `PersistentPreRunE`.
*   [x] Implement calendar logic in `internal/calendar/calendar.go` for getting the weekly events, being mindful of timezones.
*   [x] Create `cmd/week.go` to display the current week's events.
*   [x] **Milestone Verification:**
    *   [x] Build the project (`go build ./...`).
    *   [x] Request user to test the `calctl week` command.
*   [x] Implement calendar logic in `internal/calendar/calendar.go` for getting specific event details.
*   [x] Create `cmd/get.go` to get event details.
*   [x] **Milestone Verification:**
    *   [x] Build the project (`go build ./...`).
    *   [x] Request user to test the `calctl get <event-id>` command.

### Phase 3: MCP Server

*   [x] Create `cmd/mcp.go` to define the `mcp` command.
*   [x] Create `mcp/server.go` to define the MCP server.
*   [x] Implement `get_weekly_calendar` tool, reusing the logic from `internal/calendar`.
*   [x] Implement `get_event_details` tool, reusing the logic from `internal/calendar`.
*   [x] **Milestone Verification:**
    *   [x] Build the project (`go build ./...`).
    *   [x] Request user to start the MCP server and verify the tools are available.

### Phase 4: Documentation

*   [x] Create `README.md` with purpose, installation, and usage.
*   [x] Create `plans/BACKLOG.md` for future enhancements.

### Phase 5: Weekly Event Analysis

**Goal:** Add functionality to analyze the current week's events for overlaps, acceptance status, and conflicts with working hours.

*   **`internal/calendar/analysis.go`**: New file to contain the analysis logic.
*   **`cmd/week.go`**: Modify to add an `--analyze` flag.
*   **`mcp/server.go`**: Modify to add a new `analyze_weekly_calendar` tool.

#### **Tasks:**

1.  **Fetch Working Hours:**
    *   [ ] Implement logic in `internal/calendar/calendar.go` to get the user's working hours from their calendar settings using the `calendar.Settings.Get("workingHours").Do()` call.
    *   [ ] Create a new command `cmd/work-hours.go` to display the user's configured working hours.
    *   [ ] **Milestone Verification:**
        *   [ ] Build the project (`go build ./...`).
        *   [ ] Request user to test the `calctl work-hours` command.

2.  **Implement Core Analysis Logic:**
    *   [ ] Create a new file `internal/calendar/analysis.go`.
    *   [ ] In this file, create a function `AnalyzeEvents` that takes a list of events and the user's working hours.
    *   [ ] This function will perform the analysis and identify:
        *   Any events that overlap with each other.
        *   The user's response status (`accepted`, `tentative`, `needsAction`) for each event.
        *   Any events that fall outside of the user's defined working hours.
    *   [ ] The function will return a structured analysis report.

3.  **Update CLI Command:**
    *   [ ] Add a boolean `--analyze` flag to the `cmd/week.go` command.
    *   [ ] When this flag is used, the command will call the new `AnalyzeEvents` function and print a formatted report to the console.
    *   [ ] **Milestone Verification:**
        *   [ ] Build the project (`go build ./...`).
        *   [ ] Request user to test the `calctl week --analyze` command.

4.  **Create New MCP Tool:**
    *   [ ] Add a new tool named `analyze_weekly_calendar` to `mcp/server.go`.
    *   [ ] The handler for this tool will call the same core `AnalyzeEvents` logic.
    *   [ ] The tool will return the structured analysis report.
    *   [ ] **Milestone Verification:**
        *   [ ] Build the project (`go build ./...`).
        *   [ ] Request user to start the MCP server and verify the `analyze_weekly_calendar` tool is available.