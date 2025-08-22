package main

import (
	"main/crypto"
	"main/db/DAO/postgres"
	"main/log/proxies"
	"main/network/handlers"
	"main/network/middleware"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	host := os.Getenv("HOST")
	if host == "" {
		panic("following variables is not specified in env: NUM_OF_DB_CONNECTIONS")
	}

	tm := crypto.CreateTokenManager()
	daoWithLog := proxies.CreateDAOWithLog(postgres.CreateMarketplaceDAO())
	defer daoWithLog.Close()
	h := handlers.CreateHandlers(tm, daoWithLog)

	routerPaths := mux.NewRouter()
	routerPaths.HandleFunc("/login", h.Login)
	routerPaths.HandleFunc("/register", h.Register)
	routerPaths.HandleFunc("/new_announcement", h.NewAnnouncement)
	routerPaths.HandleFunc("/announcements", h.Announcements)

	authMiddleware := middleware.CreateAuthorizationMiddleware(tm)
	authMiddleware.SetNext(routerPaths)

	http.ListenAndServe(host, authMiddleware)
}
