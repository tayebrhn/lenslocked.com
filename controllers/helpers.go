package controllers

import (
	"github.com/gorilla/schema"
	"net/http"
)

func parseForm(req *http.Request, dst interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, req.PostForm); err != nil {
		return err
	}

	return nil
}
