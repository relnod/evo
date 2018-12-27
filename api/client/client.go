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
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/evo"
	"github.com/relnod/evo/pkg/math64"
)

// Client implements evo.Producer
type Client struct {
	addr    string
	conn    net.Conn
	decoder *json.Decoder

	shouldClose bool

	entitiesChangedSubscriptions map[uuid.UUID]evo.EntitiesChangedFn
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

		entitiesChangedSubscriptions: make(map[uuid.UUID]evo.EntitiesChangedFn),
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
		case api.EventCreatures:
			var creatures []*entity.Creature
			err = json.Unmarshal(event.Message, &creatures)
			if err != nil {
				log.Printf("Failed to decode creatures (%s)", err)
				return err
			}
			for _, fn := range c.entitiesChangedSubscriptions {
				fn(creatures)
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

// Size retrieves the size from the remote simulation.
func (c *Client) Size() (int, int, error) {
	resp, err := http.Get("http://" + c.addr + "/size")
	if err != nil {
		return 0, 0, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	resp.Body.Close()
	var v math64.Vec2
	err = json.Unmarshal(data, &v)
	if err != nil {
		return 0, 0, err
	}

	return int(v.X), int(v.Y), nil
}

// Creatures retrieves the creatures from the remote simulation.
func (c *Client) Creatures() ([]*entity.Creature, error) {
	resp, err := http.Get("http://" + c.addr + "/creatures")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	var creatures []*entity.Creature
	err = json.Unmarshal(data, creatures)
	if err != nil {
		return nil, err
	}

	return creatures, nil
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

func (c *Client) SubscribeEntitiesChanged(fn evo.EntitiesChangedFn) uuid.UUID {
	u := uuid.New()
	c.entitiesChangedSubscriptions[u] = fn

	err := c.sendMessage(api.EventSubscription, api.Subscription{
		Type: api.SubscriptionCreaturesChanged,
	})
	if err != nil {
		log.Fatal(err)
	}

	return u
}

func (c *Client) UnsubscribeEntitiesChanged(id uuid.UUID) {
	// TODO: actually unsubscribe world change.
	delete(c.entitiesChangedSubscriptions, id)
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
