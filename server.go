package main

import (
	"log"
	"net/http"
)

const (
	AddSrv       = "127.0.0.1:8080"
	TemplatesDir = "."
)

func main() {
	fileSrv := http.FileServer(http.Dir(TemplatesDir))

	if err := http.ListenAndServe(AddSrv, fileSrv); err != nil {
		log.Fatalln(err)
	}
}
