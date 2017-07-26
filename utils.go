package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Render util

// Render struct
type Render struct {
	layout string
	status int
	App    App
	Data   interface{}
}

// App struct
type App struct {
	Name string
	User User
	// IsAuth bool
}

// IsAuth func
func (a *App) IsAuth() bool {
	return false
}

// HTML func
func (r *Render) HTML(w http.ResponseWriter, name string, data interface{}) {
	output, err := template.New("").Delims("{{", "}}").ParseFiles(
		fmt.Sprintf("templates/%s.html", name),
		fmt.Sprintf("templates/layouts/%s.html", r.layout),
		fmt.Sprintf("templates/partials/_nav.html"),
		// fmt.Sprintf("templates/partials/_user.html"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	templates := template.Must(output, err)
	r.Data = data
	r.App.Name = config.AppName
	// r.App.IsAuth = true
	if err := templates.ExecuteTemplate(w, r.layout, r); err != nil {
		log.Fatal(err.Error())
	}
}

// Layout method
func (r *Render) Layout(name string) *Render {
	r.layout = name
	return r
}

// Status method
func (r *Render) Status(status int) *Render {
	r.status = status
	return r
}

// SendStatus method
func (r *Render) SendStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(r.status)
}

// JSON method
func (r *Render) JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.status)
	r.status = http.StatusOK
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NewRender create new Render instance
func NewRender() *Render {
	return &Render{layout: "base", status: http.StatusOK}
}

// Form util

// Form struct
type Form struct {
	FieldNames []string
	Fields     map[string]string
	Errors     map[string]string
}

// SetFields method
func (f *Form) SetFields(fields ...string) {
	f.FieldNames = fields
	f.Fields = make(map[string]string)
	f.Errors = make(map[string]string)
	for _, fieldName := range f.FieldNames {
		f.Fields[fieldName] = ""
	}
}

// Populate method
func (f *Form) Populate(r *http.Request) {
	for _, fieldName := range f.FieldNames {
		f.Fields[fieldName] = r.FormValue(fieldName)
	}
}

// Validate method
func (f *Form) Validate() bool {
	return len(f.Errors) == 0
}
