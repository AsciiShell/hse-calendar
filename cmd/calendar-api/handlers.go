package main

import (
	"net/http"

	"github.com/asciishell/HSE_calendar/internal/storage"
	"github.com/asciishell/HSE_calendar/pkg/log"
)

type Handler struct {
	storage storage.Storage
	logger  log.Logger
	rerun   chan<- interface{}
}

func NewHandler(l log.Logger, s storage.Storage, rerun chan interface{}) *Handler {
	h := Handler{storage: s, logger: l, rerun: rerun}
	return &h
}
func (h *Handler) GetDiff(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Hi", http.StatusNotFound)
}

func (h *Handler) Rerun(w http.ResponseWriter, r *http.Request) {
	h.rerun <- nil
	http.Error(w, "Task created", http.StatusOK)
}
