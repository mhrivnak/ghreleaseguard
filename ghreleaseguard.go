package main

import (
	"github.com/mhrivnak/ghreleaseguard/config"
	"github.com/mhrivnak/ghreleaseguard/push"
	"log"
	"net/http"
)

func main() {
	config.LoadConfig()
	http.HandleFunc("/api/v1/push", push.Handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
