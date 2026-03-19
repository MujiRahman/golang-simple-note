package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MujiRahman/golang-simple-note/internal/service"
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

func (c *NoteController) Create(ctx *gin.Context) {
	userID := ctx.GetUint(string(contextkey.UserIDKey))
	var req createNoteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	n, err := c.noteSvc.Create(userID, req.Title, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, n)
}

func (c *NoteController) List(ctx *gin.Context) {
	userID := ctx.GetUint(string(contextkey.UserIDKey))
	notes, err := c.noteSvc.ListByUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, notes)
}

func (c *NoteController) Get(ctx *gin.Context) {
	userID := ctx.GetUint(string(contextkey.UserIDKey))
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	n, err := c.noteSvc.GetByID(userID, uint(id64))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, n)
}

func (c *NoteController) Update(ctx *gin.Context) {
	userID := ctx.GetUint(string(contextkey.UserIDKey))
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req createNoteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	n, err := c.noteSvc.Update(userID, uint(id64), req.Title, req.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, n)
}

func (c *NoteController) Delete(ctx *gin.Context) {
	userID := ctx.GetUint(string(contextkey.UserIDKey))
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.noteSvc.Delete(userID, uint(id64)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
