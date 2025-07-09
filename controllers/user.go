package controllers
import (
  "fmt"
  "net/http"
  "lenslocked.com/views"
  "lenslocked.com/models"
)

func NewUser(us *models.UserService) *User {
  return &User {
    newView: views.NewView("bootstrap","user/new"),
    us: us,
  }
}

type User struct {
  newView *views.View
  us *models.UserService
}

func (u *User) New(res http.ResponseWriter, req *http.Request) {
  err := u.newView.Render(res,nil)
  if err != nil {
    panic(err)
  }
}

func (u *User) Create(res http.ResponseWriter, req *http.Request) {
  var form SignUpForm
  if err := parseForm(req,&form); err != nil {
    panic(err)
  }

  user := models.User{
    Name: form.Name,
    Email: form.Email,
  }

  if err := u.us.Create(&user); err != nil {
    http.Error(res, err.Error(), http.StatusInternalServerError)
    return
  }
  fmt.Fprintln(res,"User is: ",user)
}

type SignUpForm struct {
  Name string `schema:"name"`
  Email string `schema:"email"`
  Password string `schema:"password"`
}
