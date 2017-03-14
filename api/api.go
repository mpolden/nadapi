package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mpolden/nadapi/nad"
)

// API represents an API server.
type API struct {
	Client    *nad.Client
	StaticDir string
	cache     map[string]nad.Reply
	mu        sync.RWMutex
}

// Error represents an error in the API, which is returned to the user.
type Error struct {
	err     error
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (a *API) cacheSet(r nad.Reply) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.cache[r.Variable] = r
}

func (a *API) cacheGet(k string) (nad.Reply, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	v, ok := a.cache[k]
	return v, ok
}

// DeviceHandler is the handler which handles communication with an amplifier.
func (a *API) DeviceHandler(w http.ResponseWriter, req *http.Request) (interface{}, *Error) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	var cmd nad.Cmd
	if err := decoder.Decode(&cmd); err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON",
		}
	}
	reply, err := a.Client.SendCmd(cmd)
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "Failed to send command to amplifier",
		}
	}
	// Update cached value
	a.cacheSet(reply)
	return reply, nil
}

// StateHandler handles queries for the amplifiers state
func (a *API) StateHandler(w http.ResponseWriter, req *http.Request) (interface{}, *Error) {
	vars := mux.Vars(req)
	v, ok := vars["variable"]
	if !ok {
		return nil, &Error{
			Status:  http.StatusBadRequest,
			Message: "Missing required parameter: variable",
		}
	}
	_, refresh := req.URL.Query()["refresh"]
	// If not forcing refresh, return cached value if it exists
	if !refresh {
		if reply, ok := a.cacheGet(v); ok {
			return reply, nil
		}
	}
	// Send command and cache result
	reply, err := a.Client.SendCmd(nad.Cmd{Variable: v, Operator: "?"})
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "Failed to send command to amplifier",
		}
	}
	a.cacheSet(reply)
	return reply, nil
}

// NotFoundHandler handles requests to invalid routes.
func (a *API) NotFoundHandler(w http.ResponseWriter, req *http.Request) (interface{}, *Error) {
	return nil, &Error{
		err:     nil,
		Status:  http.StatusNotFound,
		Message: "Not found",
	}
}

// New returns an new API using client to communicate with an amplifier.
func New(client *nad.Client) *API {
	return &API{Client: client, cache: make(map[string]nad.Reply)}
}

type appHandler func(http.ResponseWriter, *http.Request) (interface{}, *Error)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, e := fn(w, r)
	if e != nil { // e is *Error, not os.Error.
		if e.err != nil {
			log.Print(e.err)
		}
		out, err := json.Marshal(e)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(e.Status)
		w.Write(out)
	} else {
		out, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		w.Write(out)
	}
}

func requestFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

// ListenAndServe listens on the TCP network address addr and starts serving the
// API.
func (a *API) ListenAndServe(addr string) error {
	r := mux.NewRouter()
	r.Handle("/api/v1/nad", appHandler(a.DeviceHandler))
	r.Handle("/api/v1/nad/state/{variable}", appHandler(a.StateHandler))
	r.NotFoundHandler = appHandler(a.NotFoundHandler)
	if a.StaticDir != "" {
		fs := http.StripPrefix("/static/", http.FileServer(http.Dir(a.StaticDir)))
		r.PathPrefix("/static/").Handler(fs)
	}
	return http.ListenAndServe(addr, requestFilter(r))
}
