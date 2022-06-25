package middlewares

import (
	"context"
	"github.com/vllvll/diploma/internal/repositories"
	"github.com/vllvll/diploma/internal/types"
	"net/http"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Auth(userRepository repositories.UserInterface, tokenRepository repositories.TokenInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("gophermart-auth-cookie")
			// Allow unauthenticated users in
			if err != nil || c == nil {
				next.ServeHTTP(rw, r)
				return
			}

			userId, err := tokenRepository.GetUserIdByToken(c.Value)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			user, err := userRepository.GetUserById(userId)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)
		})
	}
}

func ForContext(ctx context.Context) types.User {
	raw, _ := ctx.Value(userCtxKey).(types.User)

	return raw
}
