package views

import (
	// "fmt"
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	// "os"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
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

func (v *View) Render(wr http.ResponseWriter, data interface{}) {
	wr.Header().Set("Content-Type", "text/html")
	switch data.(type) {
	case Data:
	default:
		data = Data{
			Yield: data,
		}
	}
	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, data)
	if err != nil {
		http.Error(wr, "Something went wrong.", http.StatusInternalServerError)
		print("RENDER_ERROR: ",err.Error())
		return
	}
	io.Copy(wr, &buf)
}

func (v *View) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	v.Render(res, nil)
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
