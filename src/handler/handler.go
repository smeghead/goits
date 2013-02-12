package handler

import (
	"fmt"
	"net/http"
	"text/template"
)

func Tmpl(w http.ResponseWriter, templateName string, params map[string]interface{}) {
	t, _ := template.ParseFiles("template/layout.tmpl", fmt.Sprintf("template/%s.tmpl", templateName))
	t.Execute(w, params)
}

