package main

import (
	"net/http"
)

type homeHandler struct {
	render *Render
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.render.HTML(w, r, "pages/home", map[string]interface{}{"lol": "<script>alert(1)</script>"})
}

type aboutHandler struct {
	render *Render
}

func (h *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.render.HTML(w, r, "pages/about", nil)
}

type notFoundHandler struct {
	render *Render
}

func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.render.Status(http.StatusNotFound).HTML(w, r, "pages/404", nil)
}
