package handlers

import (
	"encoding/json"
	"fmt"
	"main/network/auth/tools"
	"net/http"
	"strconv"
)

func (h Handlers) Announcements(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logLabel := "Announcements():"

	if r.Method != http.MethodGet {
		http.Error(w, "expected GET", http.StatusMethodNotAllowed)
		return
	}

	user, ok := tools.CheckAuth(r.Context())
	if !ok {
		h.logger.WriteWarning(fmt.Sprintf("%s %s", logLabel, "an unauthorized user has gained access to the handler"))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	pageStr := r.URL.Query().Get("page")
	var page uint
	if pageStr != "" {
		page64, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			http.Error(w, "page requires uint", http.StatusBadRequest)
			return
		}
		page = uint(page64)
	}

	announcements, err := h.dao.GetAnnouncements(orderType, minPrice, maxPrice, page)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	for i, a := range announcements.Ans {
		if user.Login == a.AuthorLogin {
			announcements.Ans[i].Yours = true
		}
	}

	res, err := json.Marshal(announcements)
	if err != nil {
		h.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
