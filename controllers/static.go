package controllers
import (
  // "net/http"
  // "github.com/julienschmidt/httprouter"
  "lenslocked.com/views"
)

func NewStatic() *Static  {
  return &Static {
    Home: views.NewView("bootstrap","static/home"),
    Contact: views.NewView("bootstrap","static/contact"),
    FAQ: views.NewView("bootstrap","static/faq"),
  }
}

type Static struct {
  Home *views.View
  Contact *views.View
  FAQ *views.View
}
