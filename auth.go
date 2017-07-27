package main

import (
	"fmt"
	"html"
	"net/http"
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
	render *Render
}

func (l *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
					l.render.App.IsAuth = true
					l.render.App.User = &user
					if r.FormValue("rememberMe") == "on" {
						fmt.Println("on")
					}
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
	l.render.HTML(w, "auth/login", form)
}

type logoutHandler struct {
	render *Render
}

func (l *logoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("logoutHandler"))
	l.render.App.IsAuth = false
	l.render.App.User = nil
	http.Redirect(w, r, "/login", 301)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register"))
}

type meHandler struct {
	render *Render
}

func (m *meHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.render.HTML(w, "auth/me", nil)
}
