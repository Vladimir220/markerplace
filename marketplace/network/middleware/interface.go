package middleware

import "net/http"

type IMiddleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	SetNext(http.Handler)
}
