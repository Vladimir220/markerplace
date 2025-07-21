package handlers

import (
	"encoding/json"
	"main/network/auth"
	"net/http"
	"time"
)

type user struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "expected POST", http.StatusMethodNotAllowed)
		return
	}

	user := user{}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "expected json", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if len(user.Login) == 0 || len(user.Password) == 0 {
		http.Error(w, "expected fields login and password", http.StatusBadRequest)
	}

	a := auth.CreateAuthentication()
	token, err := a.Login(user.Login, user.Password)
	if err != nil {
		http.Error(w, "Неправильные username или password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth-cookie",
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(256),
	})
}
