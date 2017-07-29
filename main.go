package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var config = NewConfig()
var render = NewRender()

var scookie = securecookie.New([]byte(config.Secret), []byte(config.Secret[(len(config.Secret)/2):]))
var csrfMdw = csrf.Protect([]byte(config.Secret[(len(config.Secret)/2):]), csrf.Secure(config.Env == "production"))

func main() {
	r := mux.NewRouter()
	r.NotFoundHandler = http.Handler(&notFoundHandler{render})
	// r.Handle("/logout", authMdw(http.HandlerFunc(logoutHandler)))
	r.Handle("/", &homeHandler{render}).Methods("GET")
	r.Handle("/about", &aboutHandler{render}).Methods("GET")
	// api := r.PathPrefix("/api").Subrouter()

	r.Handle("/login", guestMdw(render, csrfMdw(&loginHandler{config, render, scookie}))).Methods("GET", "POST")
	r.Handle("/logout", authMdw(render, &logoutHandler{render})).Methods("GET")
	r.HandleFunc("/register", registerHandler).Methods("GET", "POST")
	r.Handle("/me", authMdw(render, &meHandler{render})).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(config.StaticDir))))
	r.PathPrefix("/node_modules/").Handler(http.StripPrefix("/node_modules/",
		http.FileServer(http.Dir("./node_modules"))))

	http.Handle("/", panicRecoveryMdw(cookieMdw(render, scookie, handlers.LoggingHandler(os.Stdout, r))))
	// http.Handle("/", r)

	fmt.Println("Server listening on port", config.Port)

	srv := &http.Server{
		// Handler: panicRecoveryHandler(baseMdw(handlers.LoggingHandler(os.Stdout, r))),
		Handler: context.ClearHandler(http.DefaultServeMux),
		Addr:    "127.0.0.1" + config.Port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
