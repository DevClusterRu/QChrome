package main

import (
	"QChromium/internal"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/search", internal.Search)

	log.Println("Starting webserver on 9598")
	err := http.ListenAndServe(":9598", nil)
	if err != nil {
		log.Fatal(err)
	}

}
