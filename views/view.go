package views

import (
	"bytes"
	"html/template"
	"io"
	"lenslocked.com/context"
	"log"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   = "views/layouts/"
	TemplateDir = "views/"
	TemplateExt = ".gohtml"
)

func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func (v *View) Render(wr http.ResponseWriter, req *http.Request, data interface{}) {
	var vd Data
	wr.Header().Set("Content-Type", "text/html")
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}

	vd.User = context.User(req.Context())

	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(wr, "Something went wrong.", http.StatusInternalServerError)
		log.Println("VIEW->Render()", err.Error())
		return
	}
	_, err = io.Copy(wr, &buf)
	if err != nil {
		http.Error(wr, "Something went wrong.", http.StatusInternalServerError)
		return
	}
}

func (v *View) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	v.Render(res, req, nil)
}

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	temp, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Layout:   layout,
		Template: temp,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}
