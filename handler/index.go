package handler

import (
	"net/http"

	logger "github.com/alecthomas/log4go"
	"github.com/smeghead/goits/data"
)

func RegisterRoutesIndex() {
	RegisterRoute("^/$", func(w http.ResponseWriter, r *http.Request, captures []string) {
		logger.Debug("index")
		params := make(map[string]interface{})
		params["topProject"] = data.GetProject("manage")
		params["wikiContent"] = data.GetWiki("manage", "top").Content
		params["projectInfos"] = data.GetProjectInfos()

		TmplTop(w, "index", params)
		logger.Debug("index end")
	})
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
