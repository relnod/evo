package websocket

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	wsocket "github.com/goxjs/websocket"

	"github.com/relnod/evo/api"
	"github.com/relnod/evo/pkg/evo"
	"github.com/relnod/evo/pkg/world"
)

// Client implements evo.Producer
type Client struct {
	addr    string
	conn    net.Conn
	decoder *json.Decoder

	worldSubscriptions map[uuid.UUID]evo.WorldStream
}

// NewClient returns a new websocket client with a given address.
func NewClient(addr string) *Client {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}

	conn, err := wsocket.Dial(u.String(), addr)
	if err != nil {
		log.Fatal("Failed to create connection: ", err)
	}

	return &Client{
		addr:    addr,
		conn:    conn,
		decoder: json.NewDecoder(conn),

		worldSubscriptions: make(map[uuid.UUID]evo.WorldStream),
	}
}

// Start starts the client.
func (c *Client) Start() error {
	for {
		event := api.Event{}
		err := c.decoder.Decode(&event)
		if err != nil {
			log.Fatal("Failed to read:", err)
		}

		switch event.Type {
		case api.World:
			w := world.World{}
			err = json.Unmarshal(event.Message, &w)
			if err != nil {
				log.Printf("Failed to decode world (%s)", err)
				return err
			}
			for _, stream := range c.worldSubscriptions {
				stream(&w)
			}
		}
	}
}

// GetWorld returns retrieves the next world object from the server.
func (c *Client) GetWorld() (*world.World, error) {
	resp, err := http.Get("http://" + c.addr + "/world")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	var w world.World
	err = json.Unmarshal(data, &w)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

// RegisterStream registers a stream via the websocket connection.
func (c *Client) SubscribeWorld(stream evo.WorldStream) uuid.UUID {
	u := uuid.New()
	c.worldSubscriptions[u] = stream

	return u
}

// UnRegisterStream un registers a stream.
func (c *Client) UnsubscribeWorld(id uuid.UUID) {
	delete(c.worldSubscriptions, id)
}
