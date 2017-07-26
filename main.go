package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var config = NewConfig()
var render = NewRender()

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/about", aboutHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/register", registerHandler)
	r.HandleFunc("/me", meHandler)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(config.StaticDir))))
	r.PathPrefix("/node_modules/").Handler(http.StripPrefix("/node_modules/",
		http.FileServer(http.Dir("./node_modules"))))

	fmt.Println("Server listening on port", config.Port)
	srv := &http.Server{
		Handler: panicRecoveryHandler(handlers.LoggingHandler(os.Stdout, r)),
		Addr:    "127.0.0.1" + config.Port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	render.HTML(w, "pages/home", nil)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	render.HTML(w, "pages/about", nil)
}
