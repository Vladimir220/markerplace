package handlers

import (
	"database/sql"
	"marketplace/network/auth/tools"
	"net/http"
	"strconv"
)

func (h Handlers) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodGet {
		http.Error(w, "expected GET", http.StatusMethodNotAllowed)
		return
	}

	user, ok := tools.CheckAuth(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	announcementIdStr := r.URL.Query().Get("announcement_id")
	var announcementId uint
	if announcementIdStr != "" {
		announcementId64, err := strconv.ParseUint(announcementIdStr, 10, 64)
		if err != nil {
			http.Error(w, "announcement_id requires uint", http.StatusBadRequest)
			return
		}
		announcementId = uint(announcementId64)
	}

	authorLogin, isAnnouncementFound, err := h.dao.GetAuthorLogin(announcementId)
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

	err = h.dao.DeleteAnnouncement(announcementId)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
