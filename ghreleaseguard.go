package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Commit struct {
	Id string
}

type Pusher struct {
	Email string
	Name  string
}

type PushMessage struct {
	Commits []Commit
	Pusher  Pusher
	Ref     string
}

func pushInspector(raw []byte) {
	var message PushMessage
	err := json.Unmarshal(raw, &message)
	if err != nil {
		log.Println("pushInspector: ", err)
		return
	}
	log.Println(message)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("pushHandler raw: ", err)
		return
	}
	go pushInspector(raw)
}

func main() {
	http.HandleFunc("/api/v1/push", pushHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
