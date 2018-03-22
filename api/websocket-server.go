package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/relnod/evo"
	"github.com/relnod/evo/num"
	"github.com/relnod/evo/world"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebsocketServer implements the internal cient interface.
type WebsocketServer struct {
	server evo.Server
	addr   string
}

// NewWebsocketServer returns a new websocket server.
func NewWebsocketServer(server evo.Server, addr string) *WebsocketServer {
	return &WebsocketServer{
		server: server,
		addr:   addr,
	}
}

// Init initializes the websocket server.
func (c *WebsocketServer) Init() {
	http.HandleFunc("/", c.handleConnection)
}

// Start starts the server.
func (c *WebsocketServer) Start() {
	go c.server.Start()

	err := http.ListenAndServe(c.addr, nil)
	if err != nil {
		log.Fatal("Failed to create server", err)
		// TODO: return err
	}
}

// handleConnection handles a websocket connection.
func (c *WebsocketServer) handleConnection(w http.ResponseWriter, r *http.Request) {

	log.Println("Connecting to Client")
	defer func() { log.Println("Disconnecting Client") }()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade to websocket failed:", err)
		return
	}
	defer conn.Close()

	id := c.server.RegisterStream(func(w *world.World) {
		conn.WriteJSON(w)
	})
	defer c.server.UnRegisterStream(id)

	for {
		event := Event{}
		err := conn.ReadJSON(&event)
		if err != nil {
			log.Println(err)
			// Disconnect
			break
		}

		switch event.Type {
		case TGetEntityAt:
			log.Println(event.Message)
			pos := &num.Vec2{}
			json.Unmarshal(event.Message, pos)
			e := c.server.GetEntityAt(pos)
			conn.WriteJSON(e)
		}
	}
}
