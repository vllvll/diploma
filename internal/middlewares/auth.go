package middlewares

import (
	"context"
	"net/http"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

type User struct {
	Name    string
	IsAdmin bool
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//c, err := r.Cookie("gophermart-auth-cookie")
		//if err != nil || c == nil {
		//	next.ServeHTTP(w, r)
		//}
		//
		//userId, err := validateAndGetUserID(c)
		//if err != nil {
		//	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		//}
		//
		//user := getUserById(db, userId)
		//ctx := context.WithValue(r.Context(), userCtxKey, user)
		//
		//r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)

	return raw
}
