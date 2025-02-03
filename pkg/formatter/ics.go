package formatter

import (
	"fmt"
	"io"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/sks/kihocche/pkg/models"
)

type ICSFormatter struct {
}

func (icsFormatter ICSFormatter) Write(writer io.Writer, events models.Events) error {
	cal := ical.NewCalendar()
	cal.SetMethod(ical.MethodPublish)
	// Convert the calendar to a string
	for i := range events {
		vEvent := icsFormatter.toICalEvent(events[i])
		cal.AddVEvent(vEvent)
	}
	_, err := writer.Write([]byte(cal.Serialize()))
	return err
}

func (ICSFormatter) toICalEvent(event models.Event) *ical.VEvent {
	vEvent := ical.NewEvent(event.ID)

	vEvent.SetStartAt(event.CreatedOn)
	vEvent.SetEndAt(event.CreatedOn.Add(time.Second))

	vEvent.SetSummary(fmt.Sprintf("[%s][%s] %s", event.Type, event.Repo.Name, event.Name))

	vEvent.SetDescription(event.Description)
	vEvent.SetOrganizer(event.Actor.Email, ical.WithCN(event.Actor.Name))
	vEvent.SetURL(event.Link)
	vEvent.SetLocation(event.Link)
	return vEvent
}
