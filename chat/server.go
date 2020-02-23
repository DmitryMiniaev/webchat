package chat

type ChatRoom struct {
	log []string
}

type Server struct {
	chatRoom *ChatRoom
	clients  map[*Client]bool
	bcast    chan string
	add      chan *Client
	remove   chan *Client
}

func NewServer() *Server {
	return &Server{
		chatRoom: &ChatRoom{log: []string{}},
		bcast:    make(chan string),
		add:      make(chan *Client),
		remove:   make(chan *Client),
		clients:  make(map[*Client]bool),
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.add:
			s.clients[client] = true
			for _, msg := range s.chatRoom.log {
				s.sendOrDetach(msg, client)
			}
		case client := <-s.remove:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
		case msg := <-s.bcast:
			s.chatRoom.log = append(s.chatRoom.log, msg)
			for client := range s.clients {
				s.sendOrDetach(msg, client)
			}
		}
	}
}

func (s *Server) sendOrDetach(msg string, client *Client) {
	select {
	case client.send <- msg:
	default:
		close(client.send)
		delete(s.clients, client)
	}
}
