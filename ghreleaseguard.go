package main

import (
	"github.com/mhrivnak/ghreleaseguard/config"
	"github.com/mhrivnak/ghreleaseguard/handlers"
	"log"
	"net/http"
)

func main() {
	config.LoadConfig()
	http.HandleFunc("/api/v1/push", handlers.PushHandler)
	http.HandleFunc("/api/v1/pullrequest", handlers.PullRequestHandler)
	err := http.ListenAndServe(config.ServerConfig.ListenAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
