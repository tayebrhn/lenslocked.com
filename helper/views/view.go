package views

import "html/template"

func NewView(layout string, files ...string) *View {

  files = append(files,
    "views/layouts/footer.gohtml",
    "views/layouts/navbar.gohtml",
    "views/layouts/bootstrap.gohtml")
  temp,err := template.ParseFiles(files...)
  if err != nil {
    panic(err)
  }

  return &View {
    Template: temp,
    Layout: layout,
  }
}

type View struct {
  Template *template.Template
  Layout string
}
