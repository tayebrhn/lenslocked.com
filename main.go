package main

import(
	"fmt"
	"net/http"
	// "github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"html/template"
)

var homeTemplate *template.Template

func home(res http.ResponseWriter,req *http.Request,_ httprouter.Params){
	res.Header().Set("Content-Type","text/html")
	// fmt.Fprint(res,"<h1>Welcome to my awesome site!</h1>")
	if err := homeTemplate.Execute(res,nil); err != nil {
		panic(err)
	}
}

func contact(res http.ResponseWriter,req *http.Request, _ httprouter.Params){
	res.Header().Set("Content-Type","text/html")
	fmt.Fprint(res,"To get in touch, please send an email"+
	"to <a href=\"mailto:support@lenslocked.com\">"+
	"support@lenslocked.com</a>",)
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

func main(){
	// r := mux.NewRouter()
	// r.HandleFunc("/",home)
	// r.HandleFunc("/contact",contact)
	// r.HandleFunc("/faq",faq)

	// var h http.Handler = http.HandlerFunc(pageNotFound)
	// r.NotFoundHandler = h
	var err error
	homeTemplate, err = template.ParseFiles("views/home.gohtml")
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	router.GET("/",home)
	router.GET("/contact",contact)
	router.GET("/faq",faq)

	// var h http.Handler = http.HandlerFunc(pageNotFound)
	// router.NotFoundHandler = h

	http.ListenAndServe(":3000",router)
}
