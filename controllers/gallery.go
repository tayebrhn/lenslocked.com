package controllers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

const (
	IndexGalleries = "index_galleries"
	EditGalleries  = "edit_galleries"
	ShowGallery    = "show_gallery"
)

const (
	maxMultipartMem = 1 << 20
)

func NewGalleries(gs models.GalleryService, router *mux.Router) *Galleries {
	return &Galleries{
		ShowView:  views.NewView("bootstrap", "gallery/show"),
		New:       views.NewView("bootstrap", "gallery/new"),
		EditView:  views.NewView("bootstrap", "gallery/edit"),
		IndexView: views.NewView("bootstrap", "gallery/index"),
		gs:        gs,
		router:    router,
	}
}

type Galleries struct {
	ShowView  *views.View
	New       *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	router    *mux.Router
}

func (g *Galleries) Index(wr http.ResponseWriter, req *http.Request) {
	user := context.User(req.Context())
	if user == nil {
		http.Error(wr, "Unauthorized", http.StatusUnauthorized)
		println(user.Name)
		return
	}
	galleries, err := g.gs.ByUserID(user.ID)

	if err != nil {
		http.Error(wr, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(wr, req, vd)
}

func (g *Galleries) Delete(wr http.ResponseWriter, req *http.Request) {
	gallery, err := g.galleriesByID(wr, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())

	if gallery.ID != user.ID {
		http.Error(wr, "You do not have permission to edit", http.StatusForbidden)
		return
	}
	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(wr, req, vd)
		return
	}

	url, err := g.router.Get(IndexGalleries).URL()
	if err != nil {
		http.Redirect(wr, req, "/", http.StatusFound)
		return
	}
	http.Redirect(wr, req, url.Path, http.StatusFound)
}

func (g *Galleries) Create(wr http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form GalleryForm

	err := parseForm(req, &form)
	if err != nil {
		vd.SetAlert(err)
		g.New.Render(wr, req, vd)
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
		g.New.Render(wr, req, vd)
		return
	}

	url, err := g.router.Get(EditGalleries).URL("id", strconv.Itoa(int(gallery.ID)))
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
	g.ShowView.Render(wr, req, vd)
}

func (g *Galleries) Edit(wr http.ResponseWriter, req *http.Request) {
	gallery, err := g.galleriesByID(wr, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())
	if gallery.UserID != user.ID {
		http.Error(wr, "You do not have permission to edit "+
			"this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(wr, req, vd)
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
		g.EditView.Render(wr, req, vd)
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
	g.EditView.Render(wr, req, vd)
}

func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleriesByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	err = r.ParseMultipartForm(maxMultipartMem)

	galleryPath := filepath.Join("images", "galleries", fmt.Sprintf("%v", gallery.ID))

	err = os.MkdirAll(galleryPath, 0755)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	files := r.MultipartForm.File["images"]

	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)

		dst, err := os.Create(filepath.Join(galleryPath, f.Filename))
		defer func(dst *os.File) {
			err := dst.Close()
			if err != nil {
				return
			}
		}(dst)
		_, err = io.Copy(dst, file)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
	}

	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Image successfully uploaded!",
	}
	g.EditView.Render(w, r, vd)
}

func (g *Galleries) galleriesByID(wr http.ResponseWriter, req *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(wr, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	fmt.Printf("INSIDE->galleriesByID g.gs: %#v\n", g.gs)
	println()
	if g.gs == nil {
		log.Fatal("g.gs is nil!")
	}
	println()

	gallery, err := g.gs.ByID(uint(id))

	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
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
