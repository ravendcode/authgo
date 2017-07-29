package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

// Render util

// Render struct
type Render struct {
	layout    string
	status    int
	App       App
	Data      interface{}
	CSRFInput template.HTML
	CSRFToken string
	// csrf.TemplateTag: csrf.TemplateField(r),
}

// App struct
type App struct {
	Name   string
	User   *User
	IsAuth bool
}

// HTML func
func (rr *Render) HTML(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	rr.SendStatus(w, rr.status)
	output, err := template.New("").Delims("{{", "}}").ParseFiles(
		fmt.Sprintf("templates/%s.html", name),
		fmt.Sprintf("templates/layouts/%s.html", rr.layout),
		fmt.Sprintf("templates/partials/_nav.html"),
		// fmt.Sprintf("templates/partials/_user.html"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	templates := template.Must(output, err)
	rr.Data = data
	rr.App.Name = config.AppName
	// rr.CSRFInput = csrf.TemplateField(r)
	token := csrf.Token(r)
	// w.Header().Set("X-CSRF-Token", token)
	rr.CSRFInput = template.HTML(fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`,
		"gorilla.csrf.Token", token))
	rr.CSRFToken = token

	if err := templates.ExecuteTemplate(w, rr.layout, rr); err != nil {
		log.Fatal(err.Error())
	}
}

// Layout method
func (rr *Render) Layout(name string) *Render {
	rr.layout = name
	return rr
}

// Status method
func (rr *Render) Status(status int) *Render {
	rr.status = status
	return rr
}

// SendStatus method
func (rr *Render) SendStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	rr.status = http.StatusOK
}

// JSON method
func (rr *Render) JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rr.SendStatus(w, rr.status)
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

// IsValid method
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}
