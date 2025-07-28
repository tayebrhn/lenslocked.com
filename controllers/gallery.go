package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/views"
)

func NewGallery(gs models.GalleryService) *Gallery {
	return &Gallery{
		NewView: views.NewView("bootstrap", "gallery/new"),
		gs: gs,
	}
}

type Gallery struct {
	NewView *views.View
	gs models.GalleryService
}

func (g *Gallery) New(wr http.ResponseWriter, req *http.Request) {
	g.NewView.Render(wr, nil)
}

func (g *Gallery) Create(wr http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form GalleryForm

	err := parseForm(req,&form)
	if err != nil {
		vd.SetAlert(err)
		g.NewView.Render(wr,vd)
		return
	}

	gallery := models.Gallery{
		Title: form.Title,
	}
	err = g.gs.Create(&gallery)
	if err != nil {
		vd.SetAlert(err)
		g.NewView.Render(wr,vd)
		return
	}
	fmt.Println(wr,gallery)
}

type GalleryForm struct {
	Title string `schema:"title"`
}