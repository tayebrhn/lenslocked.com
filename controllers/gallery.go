package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/views"
)

func NewGallery() *Gallery {
	return &Gallery{
		NewView: views.NewView("bootstrap", "gallery/new"),
	}
}

type Gallery struct {
	NewView *views.View
}

func (u *Gallery) New(res http.ResponseWriter, req *http.Request) {
	u.NewView.Render(res, nil)
}

func (u *Gallery) Create(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "// TODO: ")
}
