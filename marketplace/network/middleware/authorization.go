package middleware

import (
	"main/crypto"
	"main/network/auth"
	"net/http"
	"time"
)

const cookieFieldName = "auth-cookie"

func CreateAuthorizationMiddleware(tokenManager crypto.ITokenManager, infoLogs bool) IMiddleware {
	return &AuthorizationMiddleware{
		authorization: auth.CreateAuthorization(tokenManager, infoLogs),
	}
}

type AuthorizationMiddleware struct {
	authorization auth.IAuthorization
	next          http.Handler
}

func (am *AuthorizationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieFieldName)
	if err == http.ErrNoCookie {
		am.next.ServeHTTP(w, r)
		return
	}

	ctx, success := am.authorization.Authorize(r.Context(), cookie.Value)

	if !success {
		am.breakConnection(w)
		return
	}

	r = r.WithContext(ctx)
	am.next.ServeHTTP(w, r)
}

func (am *AuthorizationMiddleware) SetNext(next http.Handler) {
	am.next = next
}

func (am AuthorizationMiddleware) breakConnection(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieFieldName,
		Expires: time.Now().Add(-1 * time.Hour),
	})
	http.Error(w, "Authentication required", http.StatusUnauthorized)
}
