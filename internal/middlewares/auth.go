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
			if err != nil || c == nil {
				next.ServeHTTP(rw, r)
				return
			}

			userID, err := tokenRepository.GetUserIDByToken(c.Value)
			if err != nil {
				next.ServeHTTP(rw, r)
				return
			}

			user, err := userRepository.GetUserByID(userID)
			if err != nil {
				next.ServeHTTP(rw, r)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)
		})
	}
}

func ForContext(ctx context.Context) *types.User {
	value := ctx.Value(userCtxKey)
	if value == nil {
		return nil
	}

	raw, _ := value.(types.User)

	return &raw
}
