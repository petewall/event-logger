package internal

type InMemoryDatastore struct {
	events []*Event
}

func (d *InMemoryDatastore) Add(event *Event) error {
	d.events = append(d.events, event)
	return nil
}

func (d *InMemoryDatastore) List() ([]*Event, error) {
	return d.events, nil
}
