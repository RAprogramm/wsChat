// Package handlers provides application handlers.
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var (
	wsChan  = make(chan WsPayload)
	clients = make(map[WebSocketConnection]string)
	views   = jet.NewSet(
		jet.NewOSFileSystemLoader("./html"),
		jet.InDevelopmentMode(), // auto refresh template when changed
	)
	upgradeConnection = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

// Home renders the home page.
func Home(w http.ResponseWriter, _ *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

// WebSocketConnection ...
type WebSocketConnection struct {
	*websocket.Conn
}

// WsJSONResponse defines the response sent back from websocket
type WsJSONResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

// WsPayload ...
type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WsEndpoint upgrades connection to websocket
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

	go ListenForWS(&conn)
}

// ListenForWS ...
func ListenForWS(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(payload)
		if err != nil {
			log.Println(err)
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

// ListenForWSChannel ...
func ListenForWSChannel() {
	var response WsJSONResponse

	for {
		event := <-wsChan

		switch event.Action {
		case "username":
			// get a list of all users and send it back via boroadcast
			clients[event.Conn] = event.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)
		}

		// response.Action = "Got here"
		// response.Message = fmt.Sprintf("Some message, and action was %s", event.Action)
		// broadcastToAll(response)
	}
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		userList = append(userList, x)
	}
	sort.Strings(userList)
	return userList
}

func broadcastToAll(response WsJSONResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

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
