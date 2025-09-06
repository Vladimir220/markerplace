package main

import (
	"context"
	"marketplace/crypto"
	"marketplace/db/DAO/postgres"
	"marketplace/db/DAO/redis"
	"marketplace/env"
	"marketplace/log/proxies"
	"marketplace/network/handlers"
	"marketplace/network/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	host, _, err := env.GetServiceData()
	if err != nil {
		panic(err)
	}

	td, err := redis.CreateTokensDAO()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tokensDAO := proxies.CreateTokensDAOWithLog(ctx, td, true)
	tm := crypto.CreateTokenManagerProxy(ctx, tokensDAO, true)

	dao, err := postgres.CreateDAOProxy(ctx)
	if err != nil {
		panic(err)
	}
	daoWithLog := proxies.CreateDAOWithLog(ctx, dao, true)
	defer daoWithLog.Close()

	h := handlers.CreateHandlers(ctx, tm, daoWithLog, true)

	routerPaths := mux.NewRouter()
	routerPaths.HandleFunc("/login", h.Login)
	routerPaths.HandleFunc("/register", h.Register)
	routerPaths.HandleFunc("/new_announcement", h.NewAnnouncement)
	routerPaths.HandleFunc("/announcements", h.Announcements)

	authMiddleware := middleware.CreateAuthorizationMiddleware(ctx, tm, true)
	authMiddleware.SetNext(routerPaths)

	go handlers.HealthListener()

	http.ListenAndServe(host, authMiddleware)
}
