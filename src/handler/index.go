package handler

import (
	"fmt"
	"net/http"
	"./data"
)

func RegisterRoutesIndex() {
	RegisterRoute("^/$", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("xx %s", r.URL)
		params := make(map[string]interface{})
		params["projects"] = data.GetProjects()
		
		Tmpl(w, "index", params)
	})
}

