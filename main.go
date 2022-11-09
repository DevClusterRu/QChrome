package main

import (
	"QChromium/internal"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/search", internal.Search)
	http.HandleFunc("/", internal.GetImage)

	log.Println("Starting webserver on 9598")
	err := http.ListenAndServe(":9598", nil)
	if err != nil {
		log.Fatal(err)
	}

}
