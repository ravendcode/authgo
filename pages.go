package main

import "net/http"

type homeHandler struct {
	render *Render
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.render.HTML(w, "pages/home", nil)
}

type aboutHandler struct {
	render *Render
}

func (a *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.render.HTML(w, "pages/about", nil)
}
