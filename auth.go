package main

import (
	"html"
	"net/http"

	"github.com/gorilla/securecookie"
)

// User model
type User struct {
	ID       string
	Username string
	Password string
}

// Users List
type Users []User

var users = Users{
	{ID: "1", Username: "root", Password: "qwerty"},
	{ID: "2", Username: "bob", Password: "qwerty"},
	{ID: "3", Username: "vova", Password: "qwerty"},
}

type loginForm struct {
	Form
}

type loginHandler struct {
	config  *Config
	render  *Render
	scookie *securecookie.SecureCookie
}

func (h *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	form := new(loginForm)
	form.SetFields("username", "rememberMe")
	if r.Method == "POST" {
		form.Populate(r)
		if r.FormValue("username") == "" {
			form.Errors["username"] = "The username field is required."
		}
		if r.FormValue("password") == "" {
			form.Errors["password"] = "The password field is required."
		}
		if form.Validate() {
			for _, user := range users {
				if user.Username == r.FormValue("username") && user.Password == r.FormValue("password") {
					enc, err := h.scookie.Encode("auth", &user)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					cookie := &http.Cookie{
						Name:  "auth",
						Value: enc,
						Path:  "/",
					}

					if r.FormValue("rememberMe") == "on" {
						cookie.MaxAge = config.RemeberMe
					}
					http.SetCookie(w, cookie)
					if next := r.URL.Query().Get("next"); next != "" {
						http.Redirect(w, r, "/"+html.EscapeString(next), 301)
						return
					}

					http.Redirect(w, r, "/", 301)
					return
				}
			}
			form.Errors["username"] = "Incorrect username or password."
		}
	}
	h.render.HTML(w, "auth/login", form)
}

type logoutHandler struct {
	render *Render
}

func (h *logoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("logoutHandler"))
	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	h.render.App.IsAuth = false
	h.render.App.User = nil

	http.Redirect(w, r, "/login", 301)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register"))
}

type meHandler struct {
	render *Render
}

func (h *meHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.render.HTML(w, "auth/me", nil)
}
