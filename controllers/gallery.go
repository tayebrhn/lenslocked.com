package controllers

import (
	"fmt"
	"net/http"
	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

func NewGallery(gs models.GalleryService) *Gallery {
	return &Gallery{
		ShowView:views.NewView("bootstrap","gallery/show"),
		New: views.NewView("bootstrap", "gallery/new"),
		gs: gs,
	}
}

type Gallery struct {
	ShowView *views.View
	New *views.View
	gs models.GalleryService
}

func (g *Gallery) Create(wr http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form GalleryForm

	err := parseForm(req,&form)
	if err != nil {
		vd.SetAlert(err)
		g.New.Render(wr,vd)
		return
	}

	user := context.User(req.Context())

	gallery := models.Gallery{
		Title: form.Title,
		UserID: user.ID,
	}
	err = g.gs.Create(&gallery)
	if err != nil {
		vd.SetAlert(err)
		g.New.Render(wr,vd)
		return
	}
	fmt.Println(wr,gallery)
}

type GalleryForm struct {
	Title string `schema:"title"`
}