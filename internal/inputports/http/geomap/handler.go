package geomap

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Handler struct {
	templatesDir string
}

// NewHandler with TemplatesDir Constructor
func NewHandler(templatesDir string) *Handler {
	return &Handler{templatesDir: templatesDir}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join(h.templatesDir, "map.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
