package main

import (
	"html"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
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

func authMdw(render *Render, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !render.App.IsAuth {
			path := html.EscapeString(r.URL.Path)[1:]
			http.Redirect(w, r, "/login?next="+path, 302)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func cookieMdw(render *Render, scookie *securecookie.SecureCookie, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil {
			if err == http.ErrNoCookie {
				render.App.IsAuth = false
				render.App.User = nil
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var user *User
		if err = scookie.Decode("auth", cookie.Value, &user); err != nil {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			render.App.IsAuth = false
			render.App.User = nil
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userFound := false
		for _, userMemory := range users {
			if user.ID == userMemory.ID {
				userFound = true
				render.App.IsAuth = true
				render.App.User = user
				break
			}
		}
		if !userFound {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			render.App.IsAuth = false
			render.App.User = nil
		}
		next.ServeHTTP(w, r)
	})
}

func guestMdw(render *Render, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if render.App.IsAuth {
			http.Redirect(w, r, "/", 301)
			return
		}
		next.ServeHTTP(w, r)
	})
}
