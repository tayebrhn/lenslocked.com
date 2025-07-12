package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/models"
)

// func pageNotFound(res http.ResponseWriter, req *http.Request) {
// 	res.WriteHeader(http.StatusNotFound)
// 	fmt.Fprint(res, "<h1>We could not find the page you"+
// 		"were looking for:(</h1>"+
// 		"<p>Please email us if you keep being sent to an</p>"+
// 		"invalid page.</p>")
// }

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "taye"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	if err := us.DestructiveReset(); err != nil {
		panic(err)
	}

	staticController := controllers.NewStatic()
	userController := controllers.NewUser(us)
	galleryController := controllers.NewGallery()

	fileServer := http.FileServer(http.Dir("./static"))

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.Handle("/faq", staticController.FAQ).Methods("GET")
	router.HandleFunc("/signup", userController.New).Methods("GET")
	router.HandleFunc("/signup", userController.Create).Methods("POST")
	router.Handle("/login", userController.LoginView).Methods("GET")
	router.HandleFunc("/login", userController.Login).Methods("POST")
	router.HandleFunc("/gallery/new", galleryController.New).Methods("GET")
	router.HandleFunc("/cookietest", userController.CookieTest).Methods("GET")

	http.ListenAndServe(":3000", router)
}
