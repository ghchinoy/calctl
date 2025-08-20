package calendar

import (
	"fmt"
	"time"

	"google.golang.org/api/calendar/v3"
)

// GetWeeklyEvents retrieves events for the current week.
func GetWeeklyEvents(srv *calendar.Service) (*calendar.Events, error) {
	now := time.Now()
	// Sunday
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	// Saturday
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	tMin := startOfWeek.Format(time.RFC3339)
	tMax := endOfWeek.Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(tMin).TimeMax(tMax).OrderBy("startTime").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve next ten of the user's events: %v", err)
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
