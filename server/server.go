package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/caiolouro/hello-go-docker-v2/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var defaultStopTimeout = time.Second * 30

type Server struct {
	addr    string
	storage *storage.Storage
}

func NewServer(addr string, storage *storage.Storage) (*Server, error) {
	if addr == "" {
		return nil, errors.New("addr cannot be blank")
	}

	return &Server{
		addr:    addr,
		storage: storage,
	}, nil
}

// Start starts a server with a stop channel
func (s *Server) Start(stop <-chan struct{}) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.router(),
	}

	go func() {
		logrus.WithField("addr", srv.Addr).Info("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error trying to start server: %s\n", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), defaultStopTimeout)
	defer cancel()

	logrus.WithField("timeout", defaultStopTimeout).Info("stopping server")
	return srv.Shutdown(ctx)
}

func (s *Server) router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", s.defaultRoute)
	router.Methods("POST").Path("/items").Handler(Endpoint{s.createItem})
	router.Methods("GET").Path("/items").Handler(Endpoint{s.listItems})
	return router
}

func (s *Server) defaultRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}

type Endpoint struct {
	handler EndpointFunc
}

type EndpointFunc func(w http.ResponseWriter, req *http.Request) error

func (e Endpoint) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := e.handler(w, req); err != nil {
		logrus.WithError(err).Error("could not process request")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}
}
