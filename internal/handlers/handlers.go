// Package handlers provides application handlers for managing websocket connections and rendering pages.
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

// Global variables
var (
	// wsChan is a channel for websocket payloads.
	wsChan = make(chan WsPayload)

	// clients is a map that stores WebSocketConnection objects against client IDs.
	clients = make(map[WebSocketConnection]string)

	// views holds the templates used for rendering HTML pages.
	views = jet.NewSet(
		jet.NewOSFileSystemLoader("./html"),
		jet.InDevelopmentMode(),
	)

	// upgradeConnection configures the websocket upgrader with buffer sizes and an origin checker.
	upgradeConnection = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins
	}
)

// Home is an HTTP handler that renders the home page using the "home.jet" template.
func Home(w http.ResponseWriter, _ *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

// WebSocketConnection wraps a websocket.Conn in a new struct for future extensions.
type WebSocketConnection struct {
	*websocket.Conn
}

// WsJSONResponse defines the structure of messages sent over the websocket connection.
type WsJSONResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

// WsPayload defines the expected payload received from the websocket client.
type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WsEndpoint is an HTTP handler that upgrades an HTTP connection to a websocket connection.
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	var response WsJSONResponse
	response.Message = `<em><small>Connected to server</small></em>`

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}

// ListenForWs is a goroutine that listens for messages from a specific websocket connection.
func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		// Recover from panics to avoid crashing the server.
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	// Infinite loop to read messages from the websocket connection.
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// Error reading from websocket should log and break the loop.
			log.Println("Error read JSON payload")
			break
		}
		payload.Conn = *conn
		// Send the payload to the wsChan channel for processing.
		wsChan <- payload

	}
}

// ListenToWsChannel is a goroutine that processes messages sent to the wsChan channel.
func ListenToWsChannel() {
	var response WsJSONResponse

	// Infinite loop to read from the wsChan channel.
	for {
		e := <-wsChan

		// Handle different actions received from the websocket payload.
		switch e.Action {
		case "username":
			// Assign username to the client and broadcast the updated user list.
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "left":
			// Remove the client from the map and broadcast the updated user list.
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "broadcast":
			// Broadcast the received message to all connected clients.
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			broadcastToAll(response)
		}
	}
}

// getUserList compiles a sorted list of usernames of connected clients.
func getUserList() []string {
	var userList []string
	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	sort.Strings(userList)
	return userList
}

// broadcastToAll sends the provided response to all connected websocket clients.
func broadcastToAll(response WsJSONResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			// On error, log, close the connection, and remove the client from the map.
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

// renderPage renders the specified jet template with the provided data.
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
