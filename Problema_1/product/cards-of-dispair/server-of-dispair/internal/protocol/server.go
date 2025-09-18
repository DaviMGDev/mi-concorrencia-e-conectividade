package protocol

import (
	"net"
	"server-of-dispair/internal/utils"
)

type Server struct {
	Host string
	Port string

	Listener net.Listener

	Clients *utils.Map[string, *Client]

	Router *Router

	Requests  chan *Request
	Responses chan *Response
}

func NewServer(host, port string) *Server {
	s := &Server{
		Host:      host,
		Port:      port,
		Clients:   utils.NewMap[string, *Client](),
		Requests:  make(chan *Request, 100),
		Responses: make(chan *Response, 100),
	}
	s.Router = NewRouter(s)
	return s
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		return err
	}
	s.Listener = listener

	go s.Router.Start()
	go s.acceptConnections()
	go s.handleResponses()
	return nil
}

func (s *Server) Stop() error {
	for _, client := range s.Clients.Values() {
		client.Close()
	}
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			continue
		}
		client := NewClient(conn.RemoteAddr().String(), conn)
		s.Clients.Set(client.Address, client)
		go s.getRequests(client)
	}
}

func (s *Server) getRequests(client *Client) {
	for {
		request, err := client.Read()
		if err != nil {
			// On read error, assume client disconnected.
			s.Clients.Delete(client.Address)
			client.Close()
			break
		}
		s.Requests <- request
	}
}

func (s *Server) handleResponses() {
	for response := range s.Responses {
		client, exists := s.Clients.Get(response.To)
		if exists {
			client.Write(response)
		}
	}
}
