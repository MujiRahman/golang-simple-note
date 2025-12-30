package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MujiRahman/golang-simple-note/internal/helper"
	"github.com/MujiRahman/golang-simple-note/internal/service"
	"github.com/julienschmidt/httprouter"

	"github.com/MujiRahman/golang-simple-note/pkg/contextkey"
)

type NoteController struct {
	noteSvc service.NoteService
}

func NewNoteController(ns service.NoteService) *NoteController {
	return &NoteController{noteSvc: ns}
}

type createNoteReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (c *NoteController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value(contextkey.UserIDKey).(uint)
	var req createNoteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	n, err := c.noteSvc.Create(userID, req.Title, req.Content)
	if err != nil {
		helper.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	helper.RespondJSON(w, http.StatusCreated, n)
}

func (c *NoteController) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value(contextkey.UserIDKey).(uint)
	notes, err := c.noteSvc.ListByUser(userID)
	if err != nil {
		helper.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	helper.RespondJSON(w, http.StatusOK, notes)
}

func (c *NoteController) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Context().Value(contextkey.UserIDKey).(uint)
	idStr := ps.ByName("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	n, err := c.noteSvc.GetByID(userID, uint(id64))
	if err != nil {
		helper.RespondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	helper.RespondJSON(w, http.StatusOK, n)
}

func (c *NoteController) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Context().Value(contextkey.UserIDKey).(uint)
	idStr := ps.ByName("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var req createNoteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	n, err := c.noteSvc.Update(userID, uint(id64), req.Title, req.Content)
	if err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	helper.RespondJSON(w, http.StatusOK, n)
}

func (c *NoteController) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Context().Value(contextkey.UserIDKey).(uint)
	idStr := ps.ByName("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := c.noteSvc.Delete(userID, uint(id64)); err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	helper.RespondJSON(w, http.StatusNoContent, nil)
}
