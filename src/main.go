package main

import (
	"net/http"
	handler "./handler"
)

func main() {
	//static directories
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))

	http.HandleFunc("/", handler.TopHandler)
	http.ListenAndServe(":8080", nil)
}
