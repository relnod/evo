package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	wsocket "github.com/gorilla/websocket"

	"github.com/relnod/evo/api"
	"github.com/relnod/evo/pkg/evo"
	"github.com/relnod/evo/pkg/world"
)

var upgrader = wsocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Server implements evo.Consumer
type Server struct {
	producer evo.Producer
	addr     string

	subscriptions map[uuid.UUID]api.Subscription
}

// NewServer returns a new websocket server.
func NewServer(producer evo.Producer, addr string) *Server {
	s := &Server{
		producer: producer,
		addr:     addr,
	}
	http.HandleFunc("/", s.handleConnection)
	http.HandleFunc("/world", s.handleGetWorld)
	return s
}

// Start starts the server.
func (s *Server) Start() error {
	go s.producer.Start()

	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		log.Fatal("Failed to create server", err)
		return err
	}
	return nil
}

// handleConnection handles a websocket connection.
func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {

	log.Println("Connecting to Client")
	defer func() { log.Println("Disconnecting Client") }()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade to websocket failed:", err)
		return
	}
	defer conn.Close()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			// Disconnect
			break
		}
		var event api.Event
		err = json.Unmarshal(data, &event)
		if err != nil {
			log.Fatal(err)
		}
		switch event.Type {
		case api.EventSubscription:
			var subscription api.Subscription
			err = json.Unmarshal(event.Message, &subscription)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("here")
			switch subscription.Type {
			case api.SubscriptionWorld:
				id := s.producer.SubscribeWorld(func(w *world.World) {
					msg, err := json.Marshal(w)
					if err != nil {
						// TODO: maybe improve error handling
						log.Printf("Failed to unmarshal world object (%s)", err)
						return
					}
					event := &api.Event{Type: api.EventWorld, Message: msg}
					conn.WriteJSON(event)
				})
				defer s.producer.UnsubscribeWorld(id)
			}
		}
	}
}

func (s *Server) handleGetWorld(w http.ResponseWriter, r *http.Request) {
	wld, _ := s.producer.GetWorld()
	dat, err := json.Marshal(wld)
	if err != nil {
		log.Fatal(err.Error())
	}
	w.Write(dat)
}
