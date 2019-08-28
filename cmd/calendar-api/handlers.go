package main

import (
	"encoding/json"
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
	grouped, err := h.storage.GetNewLessonsFor(c)
	if err != nil {
		h.logger.Errorf("%+v:\n%+v", r, err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(grouped)
	if err != nil {
		h.logger.Errorf("%+v:\n%+v", r, err)
		return
	}
}

func (h *Handler) Rerun(w http.ResponseWriter, r *http.Request) {
	h.rerun <- nil
	http.Error(w, "Task created", http.StatusOK)
}

func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	email := chi.URLParam(r, "email")
	c := client.Client{
		Email:      email,
		GoogleCode: id,
	}
	if err := h.storage.CreateClient(&c); err != nil {
		h.logger.Errorf("%+v:\n%+v", r, err)
		http.Error(w, "can't create client", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	email := chi.URLParam(r, "email")
	c := client.Client{
		Email:      email,
		GoogleCode: id,
	}
	if err := h.storage.GetClient(&c); err != nil {
		http.Error(w, "can't find client", http.StatusNotFound)
		return
	}
	if err := h.storage.DeleteClient(&c); err != nil {
		h.logger.Errorf("%+v:\n%+v", r, err)
		http.Error(w, "can't delete client", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
