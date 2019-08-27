package main

import (
	"net/http"

	"github.com/asciishell/hse-calendar/internal/client"
	"github.com/go-chi/chi"

	"github.com/asciishell/hse-calendar/internal/storage"
	"github.com/asciishell/hse-calendar/pkg/log"
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
	id := chi.URLParam(r, "id")
	email := chi.URLParam(r, "email")
	c := client.Client{
		Email:      email,
		GoogleCode: id,
	}
	if err := h.storage.GetClient(&c); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Error(w, "Hi", http.StatusNotFound)
}

func (h *Handler) Rerun(w http.ResponseWriter, r *http.Request) {
	h.rerun <- nil
	http.Error(w, "Task created", http.StatusOK)
}
