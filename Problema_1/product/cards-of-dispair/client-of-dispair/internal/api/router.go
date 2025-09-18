package api

import (
	"client-of-dispair/internal/ui"
	"strings"
	"sync"
)

type HandlerFunc func(client *Client, chat *ui.Chat, arguments []string)

type Router struct {
	handlers map[string]HandlerFunc
	chat     *ui.Chat
	client   *Client
	quit     chan struct{}
	running  bool
	Mutex    sync.Mutex
}

func NewRouter(client *Client, chat *ui.Chat) *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
		chat:     chat,
		client:   client,
		quit:     make(chan struct{}),
	}
}

func (r *Router) AddHandler(path string, handler HandlerFunc) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.handlers[path] = handler
}

func (r *Router) Start() {
	for input := range r.chat.Inputs {
		arguments := strings.Split(input, " ")
		if !strings.HasPrefix(input, "/") {
			arguments = append([]string{"/send"}, arguments...)
		}
		path := arguments[0]
		r.Mutex.Lock()
		handler, exists := r.handlers[path]
		r.Mutex.Unlock()
		if exists {
			go handler(r.client, r.chat, arguments[1:])
		} else {
			r.chat.Outputs <- "Unknown command: " + path
		}
	}
}
