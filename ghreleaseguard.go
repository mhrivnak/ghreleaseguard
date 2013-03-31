package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func pushInspector(raw []byte) {
	var message PushMessage
	err := json.Unmarshal(raw, &message)
	if err != nil {
		log.Println("pushInspector: ", err)
		return
	}
	log.Println(message)
	log.Println(message.Release())
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
