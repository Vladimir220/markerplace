package handlers

import (
	"encoding/json"
	"fmt"
	"marketplace/network/auth/authentication"
	"marketplace/network/auth/tools"
	"net/http"
	"time"
)

func (h Handlers) Register(w http.ResponseWriter, r *http.Request) {
	logLabel := "Register():"
	defer r.Body.Close()
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
		h.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if len(user.Login) == 0 || len(user.Password) == 0 {
		http.Error(w, "expected fields login and password", http.StatusBadRequest)
	}

	a, err := authentication.CreateAuthenticationProxy(h.ctx, h.tokenManager, h.infoLogs)
	if err != nil {
		h.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	token, err := a.Register(user.Login, user.Password)
	switch err {
	case tools.ErrLoginIsTaken, tools.ErrLoginFormat, tools.ErrPasswordFormat:
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	case tools.ErrServer:
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
