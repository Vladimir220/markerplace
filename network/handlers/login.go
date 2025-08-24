package handlers

import (
	"encoding/json"
	"fmt"
	"main/network/auth"
	"net/http"
	"time"
)

type user struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h Handlers) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logLabel := "Login():"

	if r.Method != http.MethodPost {
		http.Error(w, "expected POST", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "expected json", http.StatusBadRequest)
		return
	}

	user := user{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if len(user.Login) == 0 || len(user.Password) == 0 {
		http.Error(w, "expected fields login and password", http.StatusBadRequest)
	}

	a, err := auth.CreateAuthentication(h.tokenManager, h.infoLogs)
	if err != nil {
		h.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	token, err := a.Login(user.Login, user.Password)
	switch err {
	case auth.ErrLogin:
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	case auth.ErrServer:
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth-cookie",
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 256),
	})
}
