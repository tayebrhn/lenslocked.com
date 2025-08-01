package controllers

import (
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
		gs:       gs,
		router:   router,
	}
}

type Galleries struct {
	ShowView *views.View
	New      *views.View
	gs       models.GalleryService
	router   *mux.Router
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
	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(wr, "Invalid gallery ID", http.StatusNotFound)
		return
	}

	gallery, err := g.gs.ByID(uint(id))

	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(wr, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(wr, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return
	}

	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(wr, vd)

}

type GalleryForm struct {
	Title string `schema:"title"`
}
