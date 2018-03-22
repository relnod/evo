package api

import (
	"encoding/json"
	"log"
	"net"
	"net/url"

	"github.com/goxjs/websocket"

	"github.com/relnod/evo"
	"github.com/relnod/evo/num"
	"github.com/relnod/evo/world"
	uuid "github.com/satori/go.uuid"
)

// WebsocketClient implements the internal server interface with a websocket
// connection.
type WebsocketClient struct {
	conn    net.Conn
	decoder *json.Decoder
	encoder *json.Encoder

	// TODO: only one stream should be needed here.
	streams     map[uuid.UUID]evo.Stream
	getEntityCB chan evo.GetEntityCB
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
		encoder: json.NewEncoder(conn),

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

// GetWorld retrieves the next world object from the server.
// Blocks until next world is recieved!
func (c *WebsocketClient) GetWorld() *world.World {
	w := world.World{}
	err := c.decoder.Decode(&w)
	if err != nil {
		log.Fatal("Failed to read:", err)
	}

	return &w
}

// GetEntityAt returns the entity for the given position.
// Blocks until a entity is recieved.
func (c *WebsocketClient) GetEntityAt(pos *num.Vec2, cb evo.GetEntityCB) {
	m, err := json.Marshal(&GetEntityAt{Pos: *pos})
	if err != nil {
		log.Fatal(err)
	}

	event := Event{
		Type:    TGetEntityAt,
		Message: m,
	}
	c.encoder.Encode(event)

	c.getEntityCB <- cb
	// e := entity.Creature{}
	// err = c.decoder.Decode(&e)
	// if err != nil {
	// 	log.Fatal("Failed to read:", err)
	// }
}

// RegisterStream registers a stream via the websocket connection.
// TODO: actually register the stream.
func (c *WebsocketClient) RegisterStream(stream evo.Stream) uuid.UUID {
	u := uuid.Must(uuid.NewV4())
	c.streams[u] = stream

	return u
}

// UnRegisterStream un registers a stream.
func (c *WebsocketClient) UnRegisterStream(id uuid.UUID) {
	delete(c.streams, id)
}
