package transport

import (
	"encoding/json"
	"events-exercise/internal/events"
	"log/slog"
	"net/http"
)

type EventHandler struct {
	Mux *http.ServeMux

	logger *slog.Logger
	svc    events.StreamProcessor
}

func NewHandler(logger *slog.Logger, svc events.StreamProcessor) *EventHandler {
	h := &EventHandler{
		logger: logger,
		svc:    svc,
	}

	mux := &http.ServeMux{}
	mux.Handle("/api/events", h.readEvents())
	h.Mux = mux
	return h
}

type User struct {
	Id          string                     `json:"id"`
	FullName    string                     `json:"full_name"`
	Email       string                     `json:"email"`
	BadgeCount  map[events.BadgeColour]int `json:"badge_count"`
	BadgeStatus events.BadgeStatus         `json:"badge_status"`
}

type UserResponse struct {
	Users []User `json:"users"`
}

func (h *EventHandler) outputFromResult(result map[string]*events.User) UserResponse {
	var resp UserResponse
	for _, usr := range result {
		resp.Users = append(resp.Users, User{
			Id:          usr.Id,
			FullName:    usr.FullName,
			Email:       usr.Email,
			BadgeCount:  usr.BadgeCount,
			BadgeStatus: usr.BadgeStatus(),
		})
	}
	return resp
}

func (h *EventHandler) readEvents() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(h.outputFromResult(h.svc.Result())); err != nil {
			h.logger.Error("encoding success response", "error", err)
		}
	})
}
