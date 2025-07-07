package main

import(
	"fmt"
	"net/http"
	// "github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	// "html/template"
	"lenslocked.com/views"
)

var homeView *views.View
var contactView *views.View


func home(res http.ResponseWriter,req *http.Request,_ httprouter.Params){
	res.Header().Set("Content-Type","text/html")
	errorHandle(homeView.Render(res,nil))
}

func contact(res http.ResponseWriter,req *http.Request, _ httprouter.Params){
	res.Header().Set("Content-Type","text/html")
	errorHandle(contactView.Render(res,nil))
}

func faq(res http.ResponseWriter,req *http.Request,_ httprouter.Params){
	res.Header().Set("Content-Type","text/html")
	fmt.Fprint(res,"<h1>FAQ!</h1>")
}

func pageNotFound(res http.ResponseWriter,req *http.Request,_ httprouter.Params){
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprint(res,"<h1>We could not find the page you"+
	"were looking for:(</h1>"+
	"<p>Please email us if you keep being sent to an</p>"+
	"invalid page.</p>")
}


func errorHandle(err error){
	if err != nil {
		panic(err)
	}
}

func main(){
	// r := mux.NewRouter()
	// r.HandleFunc("/",home)
	// r.HandleFunc("/contact",contact)
	// r.HandleFunc("/faq",faq)

	// var h http.Handler = http.HandlerFunc(pageNotFound)
	// r.NotFoundHandler = h
	// var err error
	homeView = views.NewView("bootstrap","views/home.gohtml")
	contactView = views.NewView("bootstrap","views/contact.gohtml")
	// homeView, err = template.ParseFiles(
	// 	"views/home.gohtml",
	// 	"views/layouts/footer.gohtml",
	// )
	// if err != nil {
	// 	panic(err)
	// }
	//
	// contactView, err = template.ParseFiles(
	// 	"views/contact.gohtml",
	// 	"views/layouts/footer.gohtml",
	// )
	// if err != nil {
	// 	panic(err)
	// }

	router := httprouter.New()
	router.GET("/",home)
	router.GET("/contact",contact)
	router.GET("/faq",faq)

	// var h http.Handler = http.HandlerFunc(pageNotFound)
	// router.NotFoundHandler = h

	http.ListenAndServe(":3000",router)
}
