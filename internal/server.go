package internal

/*
Copyright © 2024 Pete Wall <pete@petewall.net>
*/

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

const DefaultPort = 8081

type Server struct {
	Events EventList
	Port   int
}

func (s *Server) Start() error {
	log.Infof("Starting HTTP server on port %d...", s.Port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("web"))))
	r.Get("/api/events", s.getEvents)
	r.Put("/api/event", s.addEvent)
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), r)
}

func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {
	events, _ := s.Events.List()
	data, err := json.Marshal(events)
	if err != nil {
		log.WithError(err).Error("failed to convert Events into JSON")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "failed to convert Events into JSON")
		return
	}
	_, _ = w.Write(data)
}

func (s *Server) addEvent(w http.ResponseWriter, r *http.Request) {
	var event *Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.WithError(err).Error("failed to parse event")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "failed to parse event: %s", err.Error())
		return
	}

	err = s.Events.Add(event)
	if err != nil {
		log.WithError(err).Error("failed to save event")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "failed to save event")
		return
	}
}
