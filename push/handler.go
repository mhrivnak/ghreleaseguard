package push

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func inspect(raw []byte) {
	var message Message
	err := json.Unmarshal(raw, &message)
	if err != nil {
		log.Println("push.inspect: ", err)
		return
	}
	log.Println(message)
	log.Println(message.Release())
}

func Handler(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("push.Handler raw: ", err)
		return
	}
	go inspect(raw)
}
