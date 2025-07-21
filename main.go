package main

import (
	"main/handlers"
	"main/network/middleware"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		panic("following variables is not specified in env: NUM_OF_DB_CONNECTIONS")
	}

	h := handlers.CreateHandlers()
	defer h.Close()

	routerPaths := mux.NewRouter()
	routerPaths.HandleFunc("/login", h.Login)
	routerPaths.HandleFunc("/register", h.Register)
	routerPaths.HandleFunc("/new_announcement", h.NewAnnouncement)
	routerPaths.HandleFunc("/announcements", h.Announcements)

	authMiddleware := middleware.CreateAuthorizationMiddleware()
	authMiddleware.SetNext(routerPaths)

	http.ListenAndServe(host, authMiddleware)

}
