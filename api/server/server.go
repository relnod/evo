package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/relnod/evo/pkg/evo"
)

// Server implements evo.Consumer
type Server struct {
	producer evo.Producer
	addr     string

	debug bool
}

// New returns a new api server.
func New(producer evo.Producer, addr string, debug bool) *Server {
	s := &Server{
		producer: producer,
		addr:     addr,

		debug: debug,
	}

	return s
}

// Start starts the http server.
// This also starts the producer in a go routine.
func (s *Server) Start() error {
	go s.producer.Start()

	r := mux.NewRouter()
	r.HandleFunc("/connect", s.handleSocketConnection).Methods("GET")
	r.HandleFunc("/world", s.handleGetWorld).Methods("GET")
	r.HandleFunc("/stats", s.handleGetStats).Methods("GET")
	r.HandleFunc("/ticks", s.handleGetTicks).Methods("GET")
	r.HandleFunc("/ticks", s.handleSetTicks).Methods("POST")

	if s.debug {
		r.HandleFunc("/debug/pprof/", pprof.Index).Methods("GET")
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline).Methods("GET")
		r.HandleFunc("/debug/pprof/profile", pprof.Profile).Methods("GET")
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol).Methods("GET")
		r.HandleFunc("/debug/pprof/trace", pprof.Trace).Methods("GET")
	}

	err := http.ListenAndServe(s.addr, r)
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

func (s *Server) handleGetTicks(w http.ResponseWriter, r *http.Request) {
	ticks, _ := s.producer.Ticks()
	dat, err := json.Marshal(ticks)
	if err != nil {
		log.Fatal(err.Error())
	}
	w.Write(dat)
}

func (s *Server) handleSetTicks(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	ticks, err := strconv.Atoi(string(data))
	if err != nil {
		log.Fatal(err.Error())
	}
	s.producer.SetTicks(ticks)
}
