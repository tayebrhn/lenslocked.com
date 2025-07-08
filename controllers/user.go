package controllers
import (
  "fmt"
  "net/http"
  // "github.com/julienschmidt/httprouter"
  // "github.com/gorilla/schema"
  "lenslocked.com/views"
)

func NewUser() *User {
  return &User {
    NewView: views.NewView("bootstrap","views/user/new.gohtml"),
  }
}

type User struct {
  NewView *views.View
}

func (u *User) New(res http.ResponseWriter, req *http.Request) {
  err := u.NewView.Render(res,nil)
  if err != nil {
    panic(err)
  }
}

func (u *User) Create(res http.ResponseWriter, req *http.Request) {
  var form SignUpForm
  if err := parseForm(req,&form); err != nil {
    panic(err)
  }
  fmt.Fprint(res,form)
}

type SignUpForm struct {
  Email string `schema:"email"`
  Password string `schema:"password"`
}
