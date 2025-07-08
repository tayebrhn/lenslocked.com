package main

import(
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	// "github.com/julienschmidt/httprouter"
	// "html/template"
	// "lenslocked.com/views"
	"lenslocked.com/controllers"
)


// var (
// 	homeView *views.View
// 	contactView *views.View
// 	faqView *views.View
// )


// func home(res http.ResponseWriter,req *http.Request,_ httprouter.Params){
// 	res.Header().Set("Content-Type","text/html")
// 	errorHandle(homeView.Render(res,nil))
// }
//
// func contact(res http.ResponseWriter,req *http.Request, _ httprouter.Params){
// 	res.Header().Set("Content-Type","text/html")
// 	errorHandle(contactView.Render(res,nil))
// }
//
// func faq(res http.ResponseWriter,req *http.Request,_ httprouter.Params){
// 	res.Header().Set("Content-Type","text/html")
// 	errorHandle(faqView.Render(res,nil))
// }

func pageNotFound(res http.ResponseWriter,req *http.Request){
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprint(res,"<h1>We could not find the page you"+
	"were looking for:(</h1>"+
	"<p>Please email us if you keep being sent to an</p>"+
	"invalid page.</p>")
}

// func signUp(res http.ResponseWriter, req *http.Request, _ httprouter.Params)  {
// 	res.Header().Set("Content-Type","text/html")
// 	errorHandle(signUpView.Render(res,nil))
// }

// func errorHandle(err error){
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main(){

	// homeView = views.NewView("bootstrap","views/home.gohtml")
	// contactView = views.NewView("bootstrap","views/contact.gohtml")
	// faqView = views.NewView("bootstrap","views/faq.gohtml")
	staticController := controllers.NewStatic()
	userController := controllers.NewUser()
	// signUpView = views.NewView("bootstrap","views/signup.gohtml")

	fileServer := http.FileServer(http.Dir("./static"))

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",fileServer))
	router.Handle("/",staticController.Home).Methods("GET")
	router.Handle("/contact",staticController.Contact).Methods("GET")
	router.Handle("/faq",staticController.FAQ).Methods("GET")
	router.HandleFunc("/signup", userController.New).Methods("GET")
	router.HandleFunc("/signup", userController.Create).Methods("POST")

	http.ListenAndServe(":3000",router)
}
