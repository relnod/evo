package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/goxjs/websocket"

	"github.com/relnod/evo/api"
	"github.com/relnod/evo/pkg/evo"
	"github.com/relnod/evo/pkg/world"
)

// Client implements evo.Producer
type Client struct {
	addr    string
	conn    net.Conn
	decoder *json.Decoder

	shouldClose bool

	worldSubscriptions map[uuid.UUID]evo.WorldFn
}

// New returns a new websocket client with a given address.
func New(addr string) *Client {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/connect"}

	conn, err := websocket.Dial(u.String(), addr)
	if err != nil {
		log.Fatal("Failed to create connection: ", err)
	}

	return &Client{
		addr:    addr,
		conn:    conn,
		decoder: json.NewDecoder(conn),

		shouldClose: false,

		worldSubscriptions: make(map[uuid.UUID]evo.WorldFn),
	}
}

// Start starts the client.
func (c *Client) Start() error {
	for !c.shouldClose {
		event := api.Event{}
		err := c.decoder.Decode(&event)
		if err != nil {
			log.Fatal("Failed to read:", err)
		}

		switch event.Type {
		case api.EventWorld:
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

	return c.conn.Close()
}

// Stop stops the client.
func (c *Client) Stop() error {
	c.shouldClose = true
	return nil
}

// World retrieves the next world object from the server.
func (c *Client) World() (*world.World, error) {
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

// Stats retrieves the next stats object from the server.
func (c *Client) Stats() (*evo.Stats, error) {
	resp, err := http.Get("http://" + c.addr + "/stats")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	var stats evo.Stats
	err = json.Unmarshal(data, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// PauseResume toggles pause/resume of the simulation of the remote server
func (c *Client) PauseResume() error {
	_, err := http.Get("http://" + c.addr + "/pauseresume")
	return err
}

// Restart restarts the simulation on the remote server.
func (c *Client) Restart() error {
	_, err := http.Get("http://" + c.addr + "/restart")
	return err
}

// Ticks retrives the ticks per second from the server.
func (c *Client) Ticks() (int, error) {
	resp, err := http.Get("http://" + c.addr + "/ticks")
	if err != nil {
		return 0, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	ticks, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return ticks, nil
}

// SetTicks sets the ticks per second of the remote simulation.
func (c *Client) SetTicks(ticks int) error {
	_, err := http.Post("http://"+c.addr+"/ticks", "text", strings.NewReader(strconv.Itoa(ticks)))
	return err
}

func (c *Client) SubscribeWorldChange(stream evo.WorldFn) uuid.UUID {
	u := uuid.New()
	c.worldSubscriptions[u] = stream

	err := c.sendMessage(api.EventSubscription, api.Subscription{
		Type: api.SubscriptionWorld,
	})
	if err != nil {
		log.Fatal(err)
	}

	return u
}

func (c *Client) UnsubscribeWorldChange(id uuid.UUID) {
	// TODO: actually unsubscribe world change.
	delete(c.worldSubscriptions, id)
}

func (c *Client) sendMessage(t api.EventType, message interface{}) error {
	m, err := json.Marshal(message)
	if err != nil {
		return err
	}
	event := api.Event{
		Type:    t,
		Message: m,
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	c.conn.Write(data)
	return nil
}
