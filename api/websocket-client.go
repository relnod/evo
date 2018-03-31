package api

import (
	"encoding/json"
	"log"
	"net"
	"net/url"

	"github.com/goxjs/websocket"

	"github.com/relnod/evo"
	"github.com/relnod/evo/world"
	uuid "github.com/satori/go.uuid"
)

// WebsocketClient implements the internal server interface with a websocket
// connection.
type WebsocketClient struct {
	conn    net.Conn
	decoder *json.Decoder

	// TODO: only one stream should be needed here.
	streams map[uuid.UUID]evo.Stream
}

// NewWebsocketClient returns a new websocket client with a given address.
func NewWebsocketClient(addr string) *WebsocketClient {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}

	conn, err := websocket.Dial(u.String(), addr)
	if err != nil {
		log.Fatal("Failed to create connection: ", err)
	}

	return &WebsocketClient{
		conn:    conn,
		decoder: json.NewDecoder(conn),

		streams: make(map[uuid.UUID]evo.Stream),
	}
}

// Start starts the client.
func (c *WebsocketClient) Start() {
	for {
		// TODO: make it more generic.
		w := c.GetWorld()
		for _, stream := range c.streams {
			stream(w)
		}
	}
}

// GetWorld returns retrieves the next world object from the server.
// Blocks until next world is recieved!
func (c *WebsocketClient) GetWorld() *world.World {
	w := world.World{}
	err := c.decoder.Decode(&w)
	if err != nil {
		log.Fatal("Failed to read:", err)
	}

	return &w
}

// RegisterStream registers a stream via the websocket connection.
// TODO: actually register the stream.
func (c *WebsocketClient) RegisterStream(stream evo.Stream) uuid.UUID {
	u := uuid.NewV4()
	c.streams[u] = stream

	return u
}

// UnRegisterStream un registers a stream.
func (c *WebsocketClient) UnRegisterStream(id uuid.UUID) {
	delete(c.streams, id)
}
