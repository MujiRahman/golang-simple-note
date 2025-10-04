package controller

import (
	"encoding/json"
	"net/http"

	"github.com/MujiRahman/golang-simple-note/internal/helper"
	"github.com/MujiRahman/golang-simple-note/internal/service"
	"github.com/julienschmidt/httprouter"
)

type UserController struct {
	userSvc service.UserService
}

func NewUserController(us service.UserService) *UserController {
	return &UserController{userSvc: us}
}

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	u, err := c.userSvc.Register(req.Username, req.Password)
	if err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	helper.RespondJSON(w, http.StatusCreated, map[string]any{"id": u.ID, "username": u.Username})
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	token, err := c.userSvc.Login(req.Username, req.Password)
	if err != nil {
		helper.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	helper.RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}
