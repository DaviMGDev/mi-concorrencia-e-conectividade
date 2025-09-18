package protocol

import "sync"

type HandlerFunc func(server *Server, request *Request)

type Router struct {
	routes map[string]HandlerFunc
	server *Server
	mutex  sync.Mutex
}

func NewRouter(server *Server) *Router {
	return &Router{
		routes: make(map[string]HandlerFunc),
		server: server,
	}
}

func (r *Router) AddRoute(route string, handler HandlerFunc) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.routes[route] = handler
}

func (r *Router) handleRequest(request *Request) {
	r.mutex.Lock()
	handler, exists := r.routes[request.Method]
	r.mutex.Unlock()

	if exists {
		go handler(r.server, request) // Handle each request in a new goroutine
	} else {
		response := NewResponse(request.From, request.Method, "error", "Unknown method", nil)
		r.server.Responses <- response
	}
}

func (r *Router) Start() {
	for request := range r.server.Requests {
		r.handleRequest(request)
	}
}