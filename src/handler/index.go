package handler

import (
    logger "code.google.com/p/log4go"
    "net/http"
    "./data"
    "html/template"
)

func RegisterRoutesIndex() {
    RegisterRoute("^/$", func(w http.ResponseWriter, r *http.Request, captures []string) {
        logger.Debug("index")
        params := make(map[string]interface{})
        params["topProject"] = data.GetProject("manage")
        params["wikiContent"] = template.HTML(data.GetWiki("manage", "top").Content)
        params["projectInfos"] = data.GetProjectInfos()

        TmplTop(w, "index", params)
        logger.Debug("index end")
    })
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
