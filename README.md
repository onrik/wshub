# Websockets hub for golang

```go
package main

import (
    "github.com/onrik/wshub"
	"net/http"
)

var Hub = wshub.NewHub()

func WebsocketHandler(rw http.ResponseWriter, request *http.Request) {
	conn, err := Hub.NewConnection(rw, request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for message := range conn.Messages() {
		Hub.SendMessage(message)
	}
}

func main() {
	go Hub.Run()

	http.HandleFunc("/", WebsocketHandler)
	http.ListenAndServe(":8080", nil)
}

```
