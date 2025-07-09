package views

import (
  // "fmt"
  "net/http"
  "html/template"
  "path/filepath"
  // "os"
)

var (
  LayoutDir string = "views/layouts/"
  TemplateDir string = "views/"
  TemplateExt string = ".gohtml"
)

func addTemplatePath(files []string)  {
  for i,f := range files {
    files[i] = TemplateDir + f
  }
}

func addTemplateExt(files []string)  {
  for i,f := range files {
    files[i] = f + TemplateExt
  }
}

func layoutFiles() [] string {
  // cwd, err := os.Getwd()
  // if err != nil {
  //   panic(err)
  // }
  // LayoutDir := filepath.Join(cwd,"/views/layouts") + string(filepath.Separator)

  // fmt.Println("LAYOUT_DIR: ",LayoutDir)

  files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
  if err!=nil {
    panic(err)
  }
  return files
}

func NewView(layout string, files ...string) *View {
  addTemplatePath(files)
  addTemplateExt(files)
  files = append(files,layoutFiles()...)
  temp,err := template.ParseFiles(files...)
  if err != nil {
    panic(err)
  }

  return &View {
    Layout: layout,
    Template: temp,
  }
}

type View struct {
  Template *template.Template
  Layout string
}

func (v *View) Render(res http.ResponseWriter, data interface{}) error {
  res.Header().Set("Content-Type","text/html")
  return v.Template.ExecuteTemplate(res, v.Layout, data)
}

func (v *View) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
    if err := v.Render(res,nil); err != nil {
      panic(err)
    }
}
