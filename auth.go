package main

import (
	"fmt"
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	form := new(loginForm)
	form.SetFields("username", "password", "rememberMe")
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
					if r.FormValue("rememberMe") == "on" {
						fmt.Println("on")
					}
					http.Redirect(w, r, "/", 301)
					return
				}
			}
			form.Errors["username"] = "Incorrect username or password."
		}

	}
	render.HTML(w, "auth/login", form)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("logout"))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register"))
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("me"))
}
