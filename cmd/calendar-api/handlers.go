package main

import (
	"net/http"

	"github.com/asciishell/HSE_calendar/internal/storage"
	"github.com/asciishell/HSE_calendar/pkg/log"
)

type Handler struct {
	storage *storage.Storage
	logger  log.Logger
}

func NewHandler(s *storage.Storage, l log.Logger) *Handler {
	h := Handler{storage: s, logger: l}
	return &h
}
func (h *Handler) GetDiff(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Hi", http.StatusNotFound)

}
