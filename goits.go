package main

import (
	logger "code.google.com/p/log4go"
	"github.com/smeghead/goits/handler"
	"net/http"
	"runtime"
)

func main() {
	logger.LoadConfiguration("logging.xml")
	logger.Trace("main start")

	logger.Trace("CPU NUM: %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	//static directories
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))

	handler.InitRoutes()

	http.HandleFunc("/", handler.RouteHandler)
	http.ListenAndServe(":8080", nil)
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
