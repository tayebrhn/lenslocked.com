package controllers
import (
  // "net/http"
  // "github.com/julienschmidt/httprouter"
  "lenslocked.com/views"
)

func NewStatic() *Static  {
  return &Static {
    Home: views.NewView("bootstrap","views/static/home.gohtml"),
    Contact: views.NewView("bootstrap","views/static/contact.gohtml"),
    FAQ: views.NewView("bootstrap","views/static/faq.gohtml"),
  }
}

type Static struct {
  Home *views.View
  Contact *views.View
  FAQ *views.View
}
