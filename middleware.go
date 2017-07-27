package main

import (
	"html"
	"log"
	"net/http"
)

func panicRecoveryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Encountered panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func authMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if render.App.IsAuth {
			next.ServeHTTP(w, r)
			return
		}
		path := html.EscapeString(r.URL.Path)[1:]
		http.Redirect(w, r, "/login?next="+path, 301)
	})
}

func guestMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !render.App.IsAuth {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, "/me", 301)
	})
}
