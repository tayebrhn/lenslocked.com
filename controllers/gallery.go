package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

const (
	ShowGallery = "show_gallery"
)

func NewGalleries(gs models.GalleryService, router *mux.Router) *Galleries {
	return &Galleries{
		ShowView: views.NewView("bootstrap", "gallery/show"),
		New:      views.NewView("bootstrap", "gallery/new"),
		EditView: views.NewView("bootstrap", "gallery/edit"),
		gs:       gs,
		router:   router,
	}
}

type Galleries struct {
	ShowView *views.View
	New      *views.View
	EditView *views.View
	gs       models.GalleryService
	router   *mux.Router
}

func (g *Galleries) Delete(wr http.ResponseWriter, req *http.Request) {
	gallery, err := g.galleriesByID(wr, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())

	if gallery.ID != user.ID {
		http.Error(wr, "You do not have permision to edit", http.StatusForbidden)
		return
	}
	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(wr, vd)
		return
	}
	fmt.Fprintln(wr, "Succesfully deleted")
}

func (g *Galleries) Create(wr http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form GalleryForm

	err := parseForm(req, &form)
	if err != nil {
		vd.SetAlert(err)
		g.New.Render(wr, vd)
		return
	}

	user := context.User(req.Context())

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	err = g.gs.Create(&gallery)
	if err != nil {
		vd.SetAlert(err)
		g.New.Render(wr, vd)
		return
	}

	url, err := g.router.Get(ShowGallery).URL("id", strconv.Itoa(int(gallery.ID)))
	if err != nil {
		http.Redirect(wr, req, "/", http.StatusFound)
		return
	}
	http.Redirect(wr, req, url.Path, http.StatusFound)
}

func (g *Galleries) Show(wr http.ResponseWriter, req *http.Request) {
	gallery, err := g.galleriesByID(wr, req)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(wr, vd)
}

func (g *Galleries) Edit(wr http.ResponseWriter, req *http.Request) {
	gallery, err := g.galleriesByID(wr, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())
	if gallery.UserID != user.ID {
		http.Error(wr, "You do not have permision to edit "+
			"this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(wr, vd)
}

func (g *Galleries) Update(wr http.ResponseWriter, req *http.Request) {
	gallery, err := g.galleriesByID(wr, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())
	if gallery.UserID != user.ID {
		http.Error(wr, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	var form GalleryForm
	err = parseForm(req, &form)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(wr, vd)
		return
	}
	gallery.Title = form.Title

	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Gallery updated successfully",
		}
	}
	g.EditView.Render(wr, vd)
}

func (g *Galleries) galleriesByID(wr http.ResponseWriter, req *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(wr, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))

	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(wr, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(wr, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return gallery, nil
}

type GalleryForm struct {
	Title string `schema:"title"`
}
