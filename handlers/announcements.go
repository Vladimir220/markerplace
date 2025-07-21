package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h Handlers) Announcements(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	orderTypeStr := r.URL.Query().Get("order-type")
	var orderType *string
	if orderTypeStr != "" {
		orderType = &orderTypeStr
	}

	minPriceStr := r.URL.Query().Get("min_price")
	var minPrice *uint
	if minPriceStr != "" {
		minPrice64, err := strconv.ParseUint(minPriceStr, 10, 64)
		if err != nil {
			http.Error(w, "minPrice requires uint", http.StatusBadRequest)
			return
		}
		minPrice = new(uint)
		*minPrice = uint(minPrice64)
	}

	maxPriceStr := r.URL.Query().Get("max_price")
	var maxPrice *uint
	if maxPriceStr != "" {
		maxPrice64, err := strconv.ParseUint(maxPriceStr, 10, 64)
		if err != nil {
			http.Error(w, "maxPrice requires uint", http.StatusBadRequest)
			return
		}
		maxPrice = new(uint)
		*maxPrice = uint(maxPrice64)
	}

	offsetStr := r.URL.Query().Get("offset")
	var offset uint
	if offsetStr != "" {
		offset = 0
	} else {
		offset64, err := strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			http.Error(w, "offset requires uint", http.StatusBadRequest)
			return
		}
		offset = uint(offset64)
	}

	limitStr := r.URL.Query().Get("limit")
	var limit uint
	if limitStr != "" {
		limit = 10
	} else {
		limit64, err := strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			http.Error(w, "limit requires uint", http.StatusBadRequest)
			return
		}
		limit = uint(limit64)
	}

	announcement, err := h.dao.GetAnnouncements(orderType, minPrice, maxPrice, offset, limit)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(announcement)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
