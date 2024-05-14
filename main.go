package main

import (
	"fmt"
	"net/http"

	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
)

type User = map[string]interface{}
type Data = User

const port = 4000

func main() {
	// CORS
	opts := &socket.ServerOptions{}
	opts.SetCors(&types.Cors{Origin: "*"})

	users := make(map[string]User)

	io := socket.NewServer(nil, opts)
	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		fmt.Println("A user connected:", client.Id())

		client.On("user-join", func(args ...any) {
			user, ok := args[0].(User)
			if !ok {
				return
			}
			if _, ok := user["name"]; !ok {
				return
			}
			id := string(client.Id())
			fmt.Printf("User %s => %s %s joined\n", id, user["emoji"], user["name"])
			user["sid"] = id
			users[id] = user
			// Broadcast to all connected clients
			var items [][]interface{}
			for k, v := range users {
				items = append(items, []interface{}{k, v})
			}
			io.Sockets().Emit("contacts", items)
		})

		client.On("chat", func(args ...any) {
			data, ok := args[0].(Data)
			if !ok {
				return
			}
			to, ok := data["to"].(string)
			if !ok {
				return
			}
			io.Sockets().To(socket.Room(to)).Emit("chat", data)
		})

		// Create Room
		client.On("create-group", func(args ...any) {
			data, ok := args[0].(Data)
			if !ok {
				return
			}
			socketIds := data["sids"].([]interface{})
			var individualRooms []socket.Room
			for _, socketId := range socketIds {
				individualRooms = append(individualRooms, socket.Room(socketId.(string)))
			}
			roomName := data["name"].(string)
			roomId := data["id"].(string)
			// Join Room
			io.Sockets().In(individualRooms...).SocketsJoin(socket.Room(roomId))
			// Broadcast to all participants
			io.Sockets().To(socket.Room(roomId)).Emit("create-group", data)
			fmt.Printf("Room %s => %s created\n", roomId, roomName)
		})
	})

	http.Handle("/socket.io/", io.ServeHandler(nil))
	fmt.Printf("Chat server serving at localhost:%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
