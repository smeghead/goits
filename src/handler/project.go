package handler

import (
	"fmt"
	"net/http"
	"./data"
)

func ProjectTopHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("xx %s", r.URL)
	params := make(map[string]interface{})
	params["projects"] = data.GetProjects()
	
	Tmpl(w, "index", params)
}

