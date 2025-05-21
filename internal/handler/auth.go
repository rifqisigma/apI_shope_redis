package handler

import (
	"api_shope/dto"
	"api_shope/internal/usecase"
	"api_shope/utils/helper"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	jwt, err := h.authUsecase.Login(&req)
	if err != nil {
		switch err {
		case helper.ErrInvalidEmail:
			helper.WriteError(w, http.StatusBadRequest, "invalid email")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"token": jwt,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.authUsecase.Register(&req); err != nil {
		switch err {
		case helper.ErrInvalidEmail:
			helper.WriteError(w, http.StatusBadRequest, "invalid email")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, nil)
}
