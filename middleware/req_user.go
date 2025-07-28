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
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("remember_token")

		if err != nil {
			http.Redirect(wr,req,"/login",http.StatusFound)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		ctx := req.Context()
		ctx = context.WithUser(ctx,user)
		req = req.WithContext(ctx)

		if err != nil {
			http.Redirect(wr,req,"login", http.StatusFound)
			return
		}
		fmt.Println("User found: ",user)
		next(wr,req)
	})
}

func (mw *ReqUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
