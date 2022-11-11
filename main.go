package main

import (
	"QChromium/internal"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	internal.Browsers = make(map[string]*internal.Instance)

	http.HandleFunc("/search", internal.Search)
	http.HandleFunc("/", internal.GetImage)

	log.Println("Starting webserver on 9598")
	err := http.ListenAndServe(":9598", nil)
	if err != nil {
		log.Fatal(err)
	}

}
