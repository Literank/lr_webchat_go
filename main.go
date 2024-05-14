package main

import (
	"fmt"
	"net/http"

	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
)

const port = 4000

func main() {
	// CORS
	opts := &socket.ServerOptions{}
	opts.SetCors(&types.Cors{Origin: "*"})

	io := socket.NewServer(nil, opts)
	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		fmt.Println("A user connected:", client.Id())
	})

	http.Handle("/socket.io/", io.ServeHandler(nil))
	fmt.Printf("Chat server serving at localhost:%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
