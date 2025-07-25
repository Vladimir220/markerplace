package handlers

import (
	"encoding/json"
	"main/models"
	"net/http"
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
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	_, err = h.dao.NewAnnouncement(announcement)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}
