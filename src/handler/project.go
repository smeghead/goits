package handler

import (
	"fmt"
	"net/http"
	_ "./data"
)

func RegisterRoutesProject() {
	RegisterRoute("^/project", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("xx %s", r.URL)
		params := make(map[string]interface{})
		
		Tmpl(w, "project", params)
	})
}

