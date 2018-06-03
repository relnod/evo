package websocket

import (
	"encoding/json"
	"log"
	"net"
	"net/url"

	wsocket "github.com/goxjs/websocket"

	"github.com/relnod/evo"
	"github.com/relnod/evo/api"
	"github.com/relnod/evo/world"
	uuid "github.com/satori/go.uuid"
)

// Client implements the internal server interface with a websocket
// connection.
type Client struct {
	conn    net.Conn
	decoder *json.Decoder

	// TODO: only one stream should be needed here.
	streams map[uuid.UUID]evo.Stream
}

// NewClient returns a new websocket client with a given address.
func NewClient(addr string) *Client {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}

	conn, err := wsocket.Dial(u.String(), addr)
	if err != nil {
		log.Fatal("Failed to create connection: ", err)
	}

	return &Client{
		conn:    conn,
		decoder: json.NewDecoder(conn),

		streams: make(map[uuid.UUID]evo.Stream),
	}
}

// Start starts the client.
func (c *Client) Start() {
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
func (c *Client) GetWorld() *world.World {
	event := api.Event{}
	err := c.decoder.Decode(&event)
	if err != nil {
		log.Fatal("Failed to read:", err)
	}

	if event.Type != api.World {
		// TODO: actually make a request and wait for response.
		return nil
	}

	w := world.World{}
	err = json.Unmarshal(event.Message, &w)
	if err != nil {
		log.Printf("Failed to decode world (%s)", err)
	}

	return &w
}

// RegisterStream registers a stream via the websocket connection.
// TODO: actually register the stream.
func (c *Client) RegisterStream(stream evo.Stream) uuid.UUID {
	u := uuid.NewV4()
	c.streams[u] = stream

	return u
}

// UnRegisterStream un registers a stream.
func (c *Client) UnRegisterStream(id uuid.UUID) {
	delete(c.streams, id)
}
