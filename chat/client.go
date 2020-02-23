package chat

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync/atomic"
)

var clientIdGen int64

type Client struct {
	id   string
	s    *Server
	conn *websocket.Conn
	send chan string
}

func (c *Client) readWs() {
	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg := string(rawMsg)
		c.s.bcast <- newChatMsg(c.id, msg)
	}
}

func (c *Client) writeWs() {
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			payload := []byte(msg)
			if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	}
}

func (c *Client) HandleWs() {
	defer func() {
		c.s.remove <- c
		c.s.bcast <- newLeaveMsg(c.id)
		c.conn.Close()
	}()
	c.s.add <- c
	c.s.bcast <- newEnterMsg(c.id)
	go c.writeWs()
	c.readWs()
}

func NewClient(s *Server, conn *websocket.Conn) *Client {
	atomic.AddInt64(&clientIdGen, 1)
	return &Client{
		id:   strconv.FormatInt(clientIdGen, 10),
		s:    s,
		conn: conn,
		send: make(chan string, 100),
	}
}
