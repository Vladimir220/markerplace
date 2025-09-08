package handlers

import (
	"database/sql"
	"encoding/json"
	"marketplace/network/auth/tools"
	"net/http"
)

type bodyWithId struct {
	id uint `json:"id"`
}

func (h Handlers) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodDelete {
		http.Error(w, "expected DELETE", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "expected json", http.StatusBadRequest)
		return
	}

	body := bodyWithId{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.logger.WriteError("DeleteAnnouncement():" + err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	user, ok := tools.CheckAuth(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorLogin, isAnnouncementFound, err := h.dao.GetAuthorLogin(body.id)
	if err != nil {
		h.logger.WriteError("DeleteAnnouncement():" + err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if !isAnnouncementFound {
		http.Error(w, "NotFound", http.StatusNotFound)
		return
	}
	if user.Login != authorLogin && err != sql.ErrNoRows {
		if user.Group != "admin" {
			http.Error(w, "this is not your announcement", http.StatusForbidden)
			return
		}
	}

	err = h.dao.DeleteAnnouncement(body.id)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
