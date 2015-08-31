# Websockets hub for golang

Broadcasting is a most popular case for websockets. Now you can implement it in 10 code lines.

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
