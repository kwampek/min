package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterResponse struct {
	UserID    int       `json:"user_id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash)
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (A *API) createToken(user_id int, login string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user_id,
		"login":   login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(A.secret())
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
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	hash := hashPassword(req.Password)
	user_id, err := api.registerUser(req.Login, string(hash))

	if err != nil {
		http.Error(w, "user exists", http.StatusConflict)
		return
	}

	token, err := api.createToken(user_id, req.Login)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	if err := api.updateTokenBd(user_id, token); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{
		UserID:    user_id,
		Login:     req.Login,
		CreatedAt: time.Now(), // TODO
		Token:     token,
	})
}

func (api *API) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	var storedHash string
	var userId int

	err := api.DB.QueryRow(`
		SELECT user_id, password FROM messenger.users WHERE login=$1
	`, req.Login).Scan(&userId, &storedHash)

	if err == sql.ErrNoRows {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if !checkPassword(storedHash, req.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := api.createToken(userId, req.Login)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	if err := api.updateTokenBd(userId, token); err != nil {
		http.Error(w, "internal bd error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{
		UserID:    userId,
		Login:     req.Login,
		CreatedAt: time.Now(),
		Token:     token,
	})
}

func (api *API) updateTokenBd(userID int, token string) error {
	expiresAt := time.Now().Add(60 * time.Minute)

	_, err := api.DB.Exec(
		`INSERT INTO Messenger.Tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)`,
		userID,
		token,
		expiresAt,
	)

	return err
}
