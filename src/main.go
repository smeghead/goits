package main

import (
    logger "code.google.com/p/log4go"
    "net/http"
    "./handler"
)
func main() {
    logger.LoadConfiguration("logging.xml")
    logger.Trace("main start")
    //static directories
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
    http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))
    http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))

    handler.InitRoutes()

    http.HandleFunc("/", handler.RouteHandler)
    http.ListenAndServe(":8080", nil)
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
