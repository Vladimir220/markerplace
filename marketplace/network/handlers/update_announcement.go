package handlers

import (
	"encoding/json"
	"marketplace/models"
	"marketplace/network/auth/tools"
	"net/http"
	"time"
)

func (h Handlers) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
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
		h.logger.WriteError("UpdateAnnouncement():" + err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	user, ok := tools.CheckAuth(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorLogin, isAnnouncementFound, err := h.dao.GetAuthorLogin(announcement.Id)
	if err != nil {
		h.logger.WriteError("UpdateAnnouncement():" + err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if !isAnnouncementFound {
		http.Error(w, "NotFound", http.StatusNotFound)
		return
	}
	if user.Login != authorLogin {
		if user.Group != "admin" {
			http.Error(w, "this is not your announcement", http.StatusForbidden)
			return
		}
	}

	announcement.AuthorLogin = user.Login
	announcement.Date = time.Now()
	_, err = h.dao.UpdateAnnouncement(announcement)
	if err != nil {
		h.logger.WriteError("UpdateAnnouncement():" + err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}
