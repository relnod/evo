package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/relnod/evo/pkg/evo"
)

// Server implements evo.Consumer
type Server struct {
	producer evo.Producer
	addr     string
}

// New returns a new api server.
func New(producer evo.Producer, addr string) *Server {
	s := &Server{
		producer: producer,
		addr:     addr,
	}

	http.HandleFunc("/connect", s.handleSocketConnection)
	http.HandleFunc("/world", s.handleGetWorld)
	http.HandleFunc("/stats", s.handleGetStats)

	return s
}

// Start starts the http server.
// This also starts the producer in a go routine.
func (s *Server) Start() error {
	go s.producer.Start()

	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		log.Fatal("Failed to create server", err)
		return err
	}
	return nil
}

// Stop stops the server.
func (s *Server) Stop() error {
	// TODO: shutdown server
	return nil
}

func (s *Server) handleGetWorld(w http.ResponseWriter, r *http.Request) {
	wld, _ := s.producer.World()
	dat, err := json.Marshal(wld)
	if err != nil {
		log.Fatal(err.Error())
	}
	w.Write(dat)
}

func (s *Server) handleGetStats(w http.ResponseWriter, r *http.Request) {
	stats, _ := s.producer.Stats()
	dat, err := json.Marshal(stats)
	if err != nil {
		log.Fatal(err.Error())
	}
	w.Write(dat)
}
