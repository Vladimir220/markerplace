package main

import (
	"main/crypto"
	"main/db/DAO/postgres"
	"main/db/DAO/redis"
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

	td, err := redis.CreateTokensDAO()
	if err != nil {
		panic(err)
	}
	tokensDAO := proxies.CreateTokensDAOWithLog(td, true)
	tm := crypto.CreateTokenManager(tokensDAO, true)

	dao, err := postgres.CreateMarketplaceDAO()
	if err != nil {
		panic(err)
	}
	daoWithLog := proxies.CreateDAOWithLog(dao, true)
	defer daoWithLog.Close()

	h := handlers.CreateHandlers(tm, daoWithLog, true)

	routerPaths := mux.NewRouter()
	routerPaths.HandleFunc("/login", h.Login)
	routerPaths.HandleFunc("/register", h.Register)
	routerPaths.HandleFunc("/new_announcement", h.NewAnnouncement)
	routerPaths.HandleFunc("/announcements", h.Announcements)

	authMiddleware := middleware.CreateAuthorizationMiddleware(tm, true)
	authMiddleware.SetNext(routerPaths)

	http.ListenAndServe(host, authMiddleware)
}
