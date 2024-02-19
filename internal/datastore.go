package internal

import (
	"fmt"
	"strings"

	"github.com/go-logfmt/logfmt"
)

const TimeKey = "time"
const MessageKey = "msg"

type Event struct {
	Time    string            `json:"time"`
	Labels  map[string]string `json:"labels"`
	Message string            `json:"message"`
}

func (e *Event) Marshal() (string, error) {
	messageParts := []interface{}{TimeKey, e.Time}
	for key, value := range e.Labels {
		messageParts = append(messageParts, key, value)
	}
	messageParts = append(messageParts, MessageKey, e.Message)
	data, err := logfmt.MarshalKeyvals(messageParts...)
	if err != nil {
		return "", fmt.Errorf("failed to logfmt encode event: %w", err)
	}
	return string(data), nil
}

func ParseEvent(eventString string) (*Event, error) {
	event := &Event{
		Labels: map[string]string{},
	}
	d := logfmt.NewDecoder(strings.NewReader(eventString))
	for d.ScanRecord() {
		for d.ScanKeyval() {
			key := string(d.Key())
			value := string(d.Value())
			if key == TimeKey {
				event.Time = string(d.Value())
			} else if key == MessageKey {
				event.Message = value
			} else {
				event.Labels[key] = value
			}
		}
	}
	if d.Err() != nil {
		return nil, d.Err()
	}

	return event, nil
}

type EventList interface {
	Add(*Event) error
	List() ([]*Event, error)
}
