package websocket

import (
	"log"
	"net/http"

	wsocket "github.com/gorilla/websocket"
	"github.com/relnod/evo"
	"github.com/relnod/evo/world"
)

var upgrader = wsocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Server implements the internal cient interface.
type Server struct {
	server evo.Server
	addr   string
}

// NewServer returns a new websocket server.
func NewServer(server evo.Server, addr string) *Server {
	return &Server{
		server: server,
		addr:   addr,
	}
}

// Init initializes the websocket server.
func (c *Server) Init() {
	http.HandleFunc("/", c.handleConnection)
}

// Start starts the server.
func (c *Server) Start() {
	go c.server.Start()

	err := http.ListenAndServe(c.addr, nil)
	if err != nil {
		log.Fatal("Failed to create server", err)
		// TODO: return err
	}
}

// handleConnection handles a websocket connection.
func (c *Server) handleConnection(w http.ResponseWriter, r *http.Request) {

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
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Disconnect
			break
		}
	}
}
