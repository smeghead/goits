package handler

import (
    "fmt"
    "net/http"
    "text/template"
    "regexp"
)

type Route struct {
    pattern *regexp.Regexp
    h func(http.ResponseWriter, *http.Request)
}

var _routes []Route

func RegisterRoute(patternString string, handler func(http.ResponseWriter, *http.Request)) {
    fmt.Println("RegisterRoute:", patternString)
    _routes = append(_routes, Route{regexp.MustCompile(patternString), handler})
}

func InitRoutes() {
    RegisterRoutesIndex()
    RegisterRoutesProject()
    // ...
}

func RouteHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.String()

    for _, route := range _routes {
        if route.pattern.Match([]byte(path)) {
            fmt.Println("------------hit pattern:", route.pattern, "path", path)
            route.h(w, r)
            return
        }
    }
    fmt.Println("404:", path)
}

func Tmpl(w http.ResponseWriter, templateName string, params map[string]interface{}) {
    t, _ := template.ParseFiles("template/layout.tmpl", fmt.Sprintf("template/%s.tmpl", templateName))
    t.Execute(w, params)
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
