package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type FilesystemDatastore struct {
	cache *InMemoryDatastore
	Path  string
}

func (d *FilesystemDatastore) Initialize() error {
	d.cache = &InMemoryDatastore{}
	log.Debugf("Loading events from file: %s", d.Path)
	data, err := os.ReadFile(d.Path)
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			event, err := ParseEvent(line)
			if err != nil {
				return err
			}

			err = d.cache.Add(event)
			if err != nil {
				return err
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return nil
}

func (d *FilesystemDatastore) save() error {
	var lines []string
	events, _ := d.cache.List()
	for _, event := range events {
		line, err := event.Marshal()
		if err != nil {
			return err
		}
		lines = append(lines, line)
	}
	data := []byte(strings.Join(lines, "\n"))

	err := os.WriteFile(d.Path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data store file: %w", err)
	}
	return nil
}

func (d *FilesystemDatastore) Add(event *Event) error {
	err := d.cache.Add(event)
	if err != nil {
		return err
	}
	return d.save()
}

func (d *FilesystemDatastore) List() ([]*Event, error) {
	return d.cache.List()
}
