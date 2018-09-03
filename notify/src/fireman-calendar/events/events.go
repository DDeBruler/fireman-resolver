package events

import (
  "log"
  "time"

  "google.golang.org/api/calendar/v3"

  "fireman-calendar/auth"
)


func Events() *calendar.Events {
  client := calendarauth.GoogleClient()
  srv, err := calendar.New(client)
  if err != nil {
    log.Fatalf("Unable to retrieve Calendar client: %v", err)
  }

  t := time.Now().Format(time.RFC3339)
  events, err := srv.Events.List("primary").ShowDeleted(false).
    SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
  if err != nil {
    log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
  }
  return events
}
