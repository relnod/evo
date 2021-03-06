package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/relnod/evo/api"
	"github.com/relnod/evo/pkg/entity"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handleSocketConnection handles a websocket connection.
func (s *Server) handleSocketConnection(w http.ResponseWriter, r *http.Request) {
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
			switch subscription.Type {
			case api.SubscriptionCreaturesChanged:
				id := s.producer.SubscribeEntitiesChanged(func(creatures []*entity.Creature) {
					msg, err := json.Marshal(&creatures)
					if err != nil {
						// TODO: maybe improve error handling
						log.Printf("Failed to unmarshal creatures (%s)", err)
						return
					}
					event := &api.Event{Type: api.EventCreatures, Message: msg}
					conn.WriteJSON(event)
				})
				defer s.producer.UnsubscribeEntitiesChanged(id)
			}
		}
	}
}
