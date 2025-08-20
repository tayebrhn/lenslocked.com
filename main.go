package main

import (
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/middleware"
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
	router := mux.NewRouter()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer func(services *models.Services) {
		err := services.Close()
		if err != nil {
			panic(err.Error())
		}
	}(services)

	if err := services.DestructiveReset(); err != nil {
		panic(err)
	}

	//init controllers
	staticController := controllers.NewStatic()
	userController := controllers.NewUser(services.User)
	galleryController := controllers.NewGalleries(services.Gallery, router)

	//applying middlewares
	userMw := middleware.User{
		UserService: services.User,
	}
	reqUserMw := middleware.ReqUser{
		UserService: services.User,
	}

	newGallery := reqUserMw.Apply(galleryController.New)
	createGallery := reqUserMw.ApplyFn(galleryController.Create)

	//setting up static file resourses
	fileServer := http.FileServer(http.Dir("./static"))

	//routes
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.Handle("/faq", staticController.FAQ).Methods("GET")
	router.HandleFunc("/signup", userController.New).Methods("GET")
	router.HandleFunc("/signup", userController.Create).Methods("POST")
	router.Handle("/login", userController.LoginView).Methods("GET")
	router.HandleFunc("/login", userController.Login).Methods("POST")
	router.HandleFunc("/galleries", reqUserMw.ApplyFn(galleryController.Index)).Methods("GET").Name(controllers.IndexGalleries)
	router.HandleFunc("/galleries/new", newGallery).Methods("GET")
	router.HandleFunc("/galleries", createGallery).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}", galleryController.Show).Methods("GET").Name(controllers.ShowGallery)
	router.HandleFunc("/galleries/{id:[0-9]+}/edit", reqUserMw.ApplyFn(galleryController.Edit)).Methods("GET").Name(controllers.EditGalleries)
	router.HandleFunc("/galleries/{id:[0-9]+}/update", reqUserMw.ApplyFn(galleryController.Update)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/delete", reqUserMw.ApplyFn(galleryController.Delete)).Methods("POST")
	router.HandleFunc("/cookietest", userController.CookieTest).Methods("GET")

	fmt.Println("Starting server on :3000...")
	err = http.ListenAndServe(":3000", userMw.Apply(router))
	if err != nil {
		panic(err.Error())
	}
}
