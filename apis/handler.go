package apis

import (
	"net/http"
	"strings"
)

type Endpoints map[string]http.HandlerFunc

// IHandler wraps around http.Handler with some extra stuff to be used for each http handler.
type IHandler interface {
	// Prefix returns the handler's prefix, that will be ignored by the endpoints map.
	Prefix() string
	// Endpoints returns the handler's endpoints and theire {method path} combinations
	Endpoints() Endpoints
	// ServeHTTP and off course this one, that implements http.Handler
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func NewHandler(
	endpoints Endpoints,
	prefix string,
) IHandler {
	return &handler{
		endpoints: endpoints,
		prefix:    prefix,
	}
}

type handler struct {
	endpoints Endpoints
	prefix    string
}

func (h *handler) Prefix() string {
	return h.prefix
}

func (h *handler) Endpoints() Endpoints {
	return h.endpoints
}

func (h *handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	prefix := h.prefix + "/"
	endpointPath := strings.TrimPrefix(req.URL.Path, prefix[:len(prefix)-1])
	if strings.Contains(endpointPath, "/") {
		endpointPath = endpointPath[:strings.Index(endpointPath, "/")]
	}
	if handler, exists := h.Endpoints()[req.Method+" "+endpointPath]; exists {
		handler(res, req)
		return
	}
	if req.Method != http.MethodOptions {
		http.NotFound(res, req)
	}
}
