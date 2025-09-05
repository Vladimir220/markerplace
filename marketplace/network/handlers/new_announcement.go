package handlers

import (
	"encoding/json"
	"marketplace/models"
	"marketplace/network/auth/tools"
	"net/http"
	"time"
)

func (h Handlers) NewAnnouncement(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "expected POST", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "expected json", http.StatusBadRequest)
		return
	}

	announcement := models.ExtendedAnnouncement{}
	err := json.NewDecoder(r.Body).Decode(&announcement)
	if err != nil {
		h.logger.WriteError("NewAnnouncement():" + err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	user, ok := tools.CheckAuth(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	announcement.AuthorLogin = user.Login
	announcement.Date = time.Now()
	_, err = h.dao.NewAnnouncement(announcement)
	if err != nil {
		h.logger.WriteError("NewAnnouncement():" + err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}
