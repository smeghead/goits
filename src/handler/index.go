package handler

import (
    "fmt"
    "net/http"
    "./data"
)

func RegisterRoutesIndex() {
    RegisterRoute("^/$", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("index")
        params := make(map[string]interface{})
        params["topProject"] = data.GetProject("manage")
        params["wiki"] = data.GetWiki("manage", "top")
        params["projectInfos"] = data.GetProjectInfos()

        TmplTop(w, "index", params)
        fmt.Println("index end")
    })
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
