package handler

import (
    "fmt"
    "net/http"
    "net/url"
    "regexp"
    "./data"
)

func RegisterRoutesProject() {
    RegisterRoute("^/[^/]+", func(w http.ResponseWriter, r *http.Request) {
        projectNameRaw := regexp.MustCompile("^/([^/]+)").FindStringSubmatch(r.URL.RequestURI())[1]
        projectName, _ := url.QueryUnescape(projectNameRaw)
        fmt.Println("project", projectName)

        params := make(map[string]interface{})
        params["topProject"] = data.GetProject("manage")
        params["project"] = data.GetProject(projectName)

        TmplProject(w, "project", params)
    })
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
