package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	readHeaderTimeout = 60 * time.Second
)

type  KeyDataStore interface {
	Get(key string) (string, error)
}

type Server struct {
	kds KeyDataStore
	server *http.Server
	port   string
}	

func New(kds KeyDataStore, port string) (*Server, error) {
	r := mux.NewRouter()
	
	s := &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", port),
			BaseContext: func(net.Listener) context.Context {
				baseContext := context.Background()
				return baseContext
			},
			Handler:           r,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		port: port,
		kds: kds,
	}


	err := s.AddRoutes(r)
	
	return s, err
}

func (s *Server) AddRoutes(r *mux.Router) error {
	r.HandleFunc("/health", s.healthCheck).Methods(http.MethodGet)

	r = r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/value", s.getKeyVal).Methods(http.MethodGet)

	return nil
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getKeyVal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Add("Content-Type", "application/json") // TODO might do this in application specific middleware instead

	//vars := mux.Vars(r)
	vars := getBindings(r)

	keyRaw := vars["key"]
	fmt.Println(vars)
	key := keyRaw.(string)

	pN, err := s.kds.Get(key)
	if err != nil {
		fmt.Printf("\ncouldn't get value %v", err)
		handleError(ctx, w, err)
		return
	}
	handleResponse(ctx, w, pN)
}

func handleResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	jsonRes := struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}

	dataBytes, err := json.Marshal(jsonRes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("could not marshal response: %v - %v", w, err)
		return
	}

	if _, err := w.Write(dataBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("could not write response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Listen starts the server and listens on the configured port.
func (s *Server) Listen(ctx context.Context) error {
	fmt.Printf("http server starting on port: %s", s.port)

	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server error %v", err)
	}

	fmt.Println("http server stopped")

	return nil
}

func getBindings(request *http.Request) map[string]interface{} {
	// Parse the Form as part of creation if we need to
	if request.Form == nil {
		_ = request.ParseMultipartForm(32 << 20)
	}
	
	bindings := map[string]interface{}{}
	if strings.Contains(request.Header.Get("Content-Type"), "application/json") {
		body := request.Body
		defer body.Close()
		dec := json.NewDecoder(body)
		_ = dec.Decode(&bindings)

		// We should recycle putting a new body, incase other streams need to be done.
		bodyBytes, _ := json.Marshal(bindings)
		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	} else {
		key := request.URL.Query().Get("key")

		bindings["key"] = key
	}
	return bindings
}