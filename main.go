package main

import (
	"http"
	"log"
	"hajimehoshi/starrpg"
)

func main() {
	http.HandleFunc("/", starrpg.Handler)
	log.Print("Start!")
	if err := http.ListenAndServe(":3030", nil); err != nil {
		log.Fatal("http.ListenAndServe:", err)
	}
}
