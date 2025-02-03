package models

import "time"

type EventType string

const (
	EventTypePush        EventType = "push"
	EventTypePullRequest EventType = "pull_request"
	EventTypeIssue       EventType = "issue"
	EventTypeRelease     EventType = "release"
)

func (e EventType) String() string {
	return string(e)
}

func (e EventType) Color() string {
	switch e {
	case EventTypePush:
		return "green"
	case EventTypePullRequest:
		return "blue"
	case EventTypeIssue:
		return "red"
	case EventTypeRelease:
		return "purple"
	default:
		return "#f0f0f0"
	}
}

type Event struct {
	ID          string     `json:"id"`
	Type        EventType  `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Actor       EventActor `json:"actor"`
	Repo        Repo       `json:"repo"`
	Link        string     `json:"link"`
	CreatedOn   time.Time  `json:"created_on"`
}

type EventActor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// should conform to sort.Interface
type Events []Event

func (e Events) Len() int           { return len(e) }
func (e Events) Less(i, j int) bool { return e[i].CreatedOn.Before(e[j].CreatedOn) }
func (e Events) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

type EventFilter func(Event) bool

func (e Events) Add(events Events, filter ...EventFilter) Events {
	for _, event := range events {
		for _, f := range filter {
			if f(event) {
				e = append(e, event)
			}
		}
	}
	return e
}
