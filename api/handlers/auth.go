package handlers

import (
	"api/models"
	"api/service/auth"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type RegisterResponse struct {
	UserID    int       `json:"user_id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
}

func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	device := models.DeviceFromRaw(r.UserAgent(), r.RemoteAddr)

	resp, err := api.AuthService.Register(req, device)
	if err != nil {

		//http.Error(w, err.Error(), http.StatusUnauthorized)

		switch {
		case errors.Is(err, auth.ErrWeakPassword):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, auth.ErrLoginTooShort):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, auth.ErrLoginTooLong):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, auth.ErrUserAlreadyExists):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		// TODO add email and phone number check
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		log.Println("error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (api *API) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	device := models.DeviceFromRaw(r.UserAgent(), r.RemoteAddr)

	resp, err := api.AuthService.Login(req, device)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusUnauthorized)

		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, auth.ErrInvalidPassword):
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		log.Println("error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
