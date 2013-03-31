package main

import (
	"github.com/mhrivnak/ghreleaseguard/push"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/push", push.Handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
