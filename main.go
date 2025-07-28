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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()

	if err := services.DestructiveReset(); err != nil {
		panic(err)
	}

	staticController := controllers.NewStatic()
	userController := controllers.NewUser(services.User)
	galleryController := controllers.NewGallery(services.Gallery)

	reqUserMw := middleware.ReqUser{
		UserService: services.User,
	}

	newGallery := reqUserMw.Apply(galleryController.New)
	createGallery := reqUserMw.ApplyFn(galleryController.Create)

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
	router.HandleFunc("/galleries/new", newGallery).Methods("GET")
	router.HandleFunc("/galleries", createGallery).Methods("POST")
	router.HandleFunc("/cookietest", userController.CookieTest).Methods("GET")

	fmt.Printf("Starting server on :3000...")
	http.ListenAndServe(":3000", router)
}
