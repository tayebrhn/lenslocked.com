package controllers
import (
  "fmt"
  "net/http"
  "lenslocked.com/views"
  "lenslocked.com/models"
  "lenslocked.com/helpers/rand"

)

type SignUpForm struct {
  Name string `schema:"name"`
  Email string `schema:"email"`
  Age uint `schema:"age"`
  Password string `schema:"password"`
}

type loginForm struct {
  Email string `schema:"email"`
  Password string `schema:"password"`
}

type User struct {
  NewView *views.View
  LoginView *views.View
  us models.UserService
}

func NewUser(us models.UserService) *User {
  return &User {
    NewView: views.NewView("bootstrap","user/new"),
    LoginView: views.NewView("bootstrap","user/login"),
    us: us,
  }
}

func (u *User) New(wr http.ResponseWriter, req *http.Request) {
  err := u.NewView.Render(wr,nil)
  if err != nil {
    panic(err)
  }
}

func (u *User) Create(wr http.ResponseWriter, req *http.Request) {
  var form SignUpForm
  if err := parseForm(req,&form); err != nil {
    panic(err)
  }

  user := models.User{
    Name: form.Name,
    Email: form.Email,
    Age: form.Age,
    Password: form.Password,
  }

  if err := u.us.Create(&user); err != nil {
    http.Error(wr, err.Error(), http.StatusInternalServerError)
    return
  }
  err := u.signIn(wr, &user)
  if err != nil {
    http.Error(wr, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(wr, req, "/cookietest",http.StatusFound)
  // foundUser := u.us.InAgeRange(15,25)
  // fmt.Fprintln(wr,"User is: ",foundUser)
}

func (u *User) Login(wr http.ResponseWriter, req *http.Request) {
  var form loginForm
  if err := parseForm(req,&form); err != nil {
    panic(err)
  }
  user, err := u.us.Authenticate(form.Email, form.Password)

  if err != nil {
    switch err {
    case models.ErrNotFound:
      fmt.Fprintln(wr,"invalid email address")
    case models.ErrPasswordIncorrect:
      fmt.Fprintln(wr,"invalid password provided")
    default:
      http.Error(wr,err.Error(), http.StatusInternalServerError)
    }
    return
  }
  err = u.signIn(wr, user)
  if err != nil {
    http.Error(wr, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(wr, req, "/cookietest",http.StatusFound)
}

func (u *User) CookieTest(wr http.ResponseWriter, req *http.Request) {
  cookie, err := req.Cookie("remeber_token")
  if err != nil {
    http.Error(wr, err.Error(), http.StatusInternalServerError)
    return
  }
  user, err := u.us.ByRemember(cookie.Value)
  fmt.Fprintln(wr,user)
}

func (u *User) signIn(wr http.ResponseWriter, user *models.User) error  {
  if user.Remember == "" {
    token, err := rand.RememberToken()
    if err != nil {
      return err
    }
    user.Remember = token
    err = u.us.Update(user)
    if err != nil {
      return err
    }
  }
  cookie := http.Cookie{
    Name: "remeber_token",
    Value: user.Remember,
    HttpOnly: true,
  }

  http.SetCookie(wr, &cookie)
  return nil
}
