package middleware

import (
	"fmt"
	"net/http"

	"lenslocked.com/context"
	"lenslocked.com/models"
)

type ReqUser struct {
	models.UserService
}

func (mw *ReqUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		user := context.User(req.Context())
		if user == nil {
			http.Redirect(wr, req, "/login", http.StatusFound)
			return
		}
		next(wr, req)
	}
}

func (mw *ReqUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

type User struct {
	models.UserService
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("remember_token")

		if err != nil {
			next(wr, req)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(wr, req)
			return
		}

		ctx := req.Context()
		ctx = context.WithUser(ctx, user)
		req = req.WithContext(ctx)

		fmt.Println("User found: ", user)
		next(wr, req)
	}
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
