package events

import (
  "fmt"
  "fireman-calendar/events"
)

func main() {
  events := events.Events();
  fmt.Println("Upcoming events:")
  if len(events.Items) == 0 {
    fmt.Println("No upcoming events found.")
  } else {
    for _, item := range events.Items {
      date := item.Start.DateTime
      if date == "" {
        date = item.Start.Date
      }
      fmt.Printf("%v (%v)\n", item.Summary, date)
    }
  }
}
