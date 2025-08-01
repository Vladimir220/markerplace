package middleware

import (
	"main/network/auth"
	"main/tools/crypto"
	"main/tools/log"
	"net/http"
	"time"
)

const cookieFieldName = "auth-cookie"

func CreateAuthorizationMiddleware(tokenManager crypto.ITokenManager) IMiddleware {
	return &AuthorizationMiddleware{
		authorization: auth.CreateAuthorization(tokenManager),
		logger:        log.CreateLogger("AuthorizationMiddleware:"),
		tokenManager:  tokenManager,
	}
}

type AuthorizationMiddleware struct {
	authorization auth.IAuthorization
	logger        log.ILogger
	tokenManager  crypto.ITokenManager
	next          http.Handler
}

func (am *AuthorizationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieFieldName)
	if err != nil {
		am.logger.WriteInfo("ServeHTTP:" + err.Error())
		am.next.ServeHTTP(w, r)
		return
	}

	ctx, success := am.authorization.Authorize(r.Context(), cookie.Value)

	if !success {
		am.logger.WriteInfo("ServeHTTP:Unknown cookies:" + cookie.Value)
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
