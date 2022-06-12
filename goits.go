package main

import (
	"net/http"
	"os"
	"runtime"

	logger "github.com/alecthomas/log4go"
	"github.com/chai2010/gettext-go"
	"github.com/smeghead/goits/handler"
)

func main() {
	gettext.BindLocale(gettext.New("goits", "locale"))
	gettext.SetDomain("goits")
	os.Setenv("LANGUAGE", "ja_JP.utf8")
	gettext.SetLanguage(gettext.DefaultLanguage)

	logger.LoadConfiguration("logging.xml")
	logger.Trace("main start")

	logger.Trace("CPU NUM: %d", runtime.NumCPU())

	//static directories
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))

	handler.InitRoutes()

	http.HandleFunc("/", handler.RouteHandler)
	http.ListenAndServe(":8080", nil)
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
