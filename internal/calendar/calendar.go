package calendar

import (
	"fmt"
	"time"

	"google.golang.org/api/calendar/v3"
)

// GetWeeklyEvents retrieves events for the week of the given date.
// If dateStr is empty, it uses the current week.
func GetWeeklyEvents(srv *calendar.Service, dateStr string) (*calendar.Events, error) {
	var now time.Time
	var err error
	if dateStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, please use YYYY-MM-DD: %w", err)
		}
	}

	// Sunday
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	// Saturday
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	tMin := startOfWeek.Format(time.RFC3339)
	tMax := endOfWeek.Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(tMin).TimeMax(tMax).OrderBy("startTime").Fields("*").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve events for the week: %v", err)
	}

	return events, nil
}

// GetEvent retrieves a single event by its ID.
func GetEvent(srv *calendar.Service, eventID string) (*calendar.Event, error) {
	event, err := srv.Events.Get("primary", eventID).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve event %q: %v", eventID, err)
	}
	return event, nil
}

// GetWorkingHours retrieves the user's working hours settings.
func GetWorkingHours(srv *calendar.Service) (*calendar.Setting, error) {
	setting, err := srv.Settings.Get("workingHours").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve working hours setting: %v", err)
	}
	return setting, nil
}

// ListSettings retrieves all the user's settings.
func ListSettings(srv *calendar.Service) (*calendar.Settings, error) {
	settings, err := srv.Settings.List().Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve settings: %v", err)
	}
	return settings, nil
}
