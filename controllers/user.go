package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/helpers/rand"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

func (u *User) New(wr http.ResponseWriter, req *http.Request) {
	u.NewView.Render(wr, nil)
}

func (u *User) Create(wr http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form SignUpForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(wr, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Age:      form.Age,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(wr, vd)
		return
	}
	err := u.signIn(wr, &user)
	if err != nil {
		http.Redirect(wr, req, "/login", http.StatusFound)
		return
	}
	http.Redirect(wr, req, "/cookietest", http.StatusFound)
}

func (u *User) Login(wr http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form loginForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(wr, vd)
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("No user exists with that email address")
		case models.ErrPasswordIncorrect:
			vd.SetAlert(err)
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(wr, vd)
		return
	}
	err = u.signIn(wr, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(wr, vd)
		return
	}
	http.Redirect(wr, req, "/cookietest", http.StatusFound)
}

func (u *User) CookieTest(wr http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("remeber_token")
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	user, _ := u.us.ByRemember(cookie.Value)
	fmt.Fprintln(wr, user)
}

func (u *User) signIn(wr http.ResponseWriter, user *models.User) error {
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
		Name:     "remeber_token",
		Value:    user.Remember,
		HttpOnly: true,
	}

	http.SetCookie(wr, &cookie)
	return nil
}

func NewUser(us models.UserService) *User {
	return &User{
		NewView:   views.NewView("bootstrap", "user/new"),
		LoginView: views.NewView("bootstrap", "user/login"),
		us:        us,
	}
}

type User struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

type SignUpForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Age      uint   `schema:"age"`
	Password string `schema:"password"`
}

type loginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
