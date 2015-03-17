package api

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/martinp/nadapi/nad"
	"log"
	"net/http"
)

// API represents an API server.
type API struct {
	Client nad.Client
}

// Error represents an error in the API, which is returned to the user.
type Error struct {
	err     error
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func marshal(data interface{}, indent bool) ([]byte, error) {
	if indent {
		return json.MarshalIndent(data, "", "  ")
	}
	return json.Marshal(data)
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
func New(client nad.Client) API {
	return API{Client: client}
}

type appHandler func(http.ResponseWriter, *http.Request) (interface{}, *Error)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, e := fn(w, r)
	if e != nil { // e is *Error, not os.Error.
		if e.err != nil {
			log.Print(e.err)
		}
		jsonBlob, err := marshal(e, true)
		if err != nil {
			// Should never happen
			panic(err)
		}
		w.WriteHeader(e.Status)
		w.Write(jsonBlob)
	} else {
		indent := context.Get(r, "indent").(bool)
		jsonBlob, err := marshal(data, indent)
		if err != nil {
			panic(err)
		}
		w.Write(jsonBlob)
	}
}

func requestFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, indent := r.URL.Query()["pretty"]
		context.Set(r, "indent", indent)
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// ListenAndServe listens on the TCP network address addr and starts serving the
// API.
func (a *API) ListenAndServe(addr string) error {
	r := mux.NewRouter()
	r.Handle("/api/v1/nad", appHandler(a.DeviceHandler))
	r.NotFoundHandler = appHandler(a.NotFoundHandler)
	return http.ListenAndServe(addr, requestFilter(r))
}
