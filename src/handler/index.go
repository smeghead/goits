package handler

import (
	"fmt"
	"net/http"
	data "./data"
	"text/template"
)

func TopHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("xx %s", r.URL)
	params := make(map[string]interface{})
	params["projects"] = data.GetProjects()
	
	t, _ := template.ParseFiles("template/layout.tmpl", "template/index.tmpl")
	t.Execute(w, params)
}

