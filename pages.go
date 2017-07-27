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

func (h *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.render.HTML(w, "pages/about", nil)
}

type notFoundHandler struct {
	render *Render
}

func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.render.Status(http.StatusNotFound).HTML(w, "pages/404", nil)
}
