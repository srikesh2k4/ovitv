package signaling

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type      string `json:"type"`
	SDP       string `json:"sdp,omitempty"`
	Candidate string `json:"candidate,omitempty"`
	Text      string `json:"text,omitempty"`
	IP        string `json:"ip,omitempty"`
	Gender    string `json:"gender,omitempty"`
	Age       int    `json:"age,omitempty"`
	Interest  string `json:"interest,omitempty"`
	MinAge    int    `json:"minAge,omitempty"`
	MaxAge    int    `json:"maxAge,omitempty"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte

	gender   string
	age      int
	interest string
	minAge   int
	maxAge   int
}

type Hub struct {
	clients    map[*Client]bool
	partners   map[*Client]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		partners:   make(map[*Client]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)

		case msg := <-h.broadcast:
			h.handleBroadcast(msg)
		}
	}
}

// ---------------- HANDLE REGISTER ----------------

func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true

	// Find partner
	for c := range h.clients {
		if c == client || h.partners[c] != nil {
			continue
		}

		if (client.interest == "any" || c.gender == client.interest) &&
			(c.interest == "any" || client.gender == c.interest) &&
			c.age >= client.minAge && c.age <= client.maxAge &&
			client.age >= c.minAge && client.age <= c.maxAge {

			h.partners[client] = c
			h.partners[c] = client

			// send match notification: send a's info to b, and b's info to a
			sendMatch(client, c)
			sendMatch(c, client)
			break
		}
	}
}

// ---------------- HANDLE UNREGISTER ----------------

func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, client)

	if p := h.partners[client]; p != nil {
		// notify partner that the other left (optional)
		notify := Message{Type: "partner_left"}
		if data, err := json.Marshal(notify); err == nil {
			select {
			case p.send <- data:
			default:
				// best-effort: close and remove
				close(p.send)
				delete(h.clients, p)
			}
		}

		delete(h.partners, p)
		delete(h.partners, client)
	}
}

// ---------------- HANDLE BROADCAST ----------------

func (h *Hub) handleBroadcast(msg []byte) {
	var m Message
	if json.Unmarshal(msg, &m) != nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.clients {
		if p := h.partners[c]; p != nil {
			// send to partner only
			select {
			case p.send <- msg:
			default:
				// partner's send channel is full - drop and clean up
				close(p.send)
				delete(h.clients, p)
				if pp := h.partners[p]; pp != nil {
					delete(h.partners, pp)
				}
				delete(h.partners, p)
			}
		}
	}
}

func sendMatch(a, b *Client) {
	// send a's public info (IP) to b
	ip := ""
	if a.conn != nil {
		ip = a.conn.RemoteAddr().String()
	}
	data, _ := json.Marshal(Message{
		Type: "match",
		IP:   ip,
	})
	// send to b (the receiver)
	select {
	case b.send <- data:
	default:
		// if b's send is blocked, close and cleanup
		close(b.send)
	}
}

// ---------------- WEBSOCKET HANDLER ----------------

// ServeWs returns an http.Handler compatible with middleware
func ServeWs(hub *Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte, 256),
		}

		// start pumps
		go client.writePump()
		go client.readPump()
	})
}

// ---------------- CLIENT READ PUMP ----------------

const (
	pongWait       = 60 * time.Second
	pingPeriod     = 30 * time.Second
	maxMessageSize = 64 << 10 // 64KB
)

func (c *Client) readPump() {
	defer func() {
		// ensure unregister and close
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			// websocket closed or error
			break
		}

		var m Message
		if json.Unmarshal(msg, &m) != nil {
			continue
		}

		if m.Type == "join" {
			c.gender = m.Gender
			c.age = m.Age
			c.interest = m.Interest
			c.minAge = m.MinAge
			c.maxAge = m.MaxAge
			c.hub.register <- c
		} else {
			c.hub.broadcast <- msg
		}
	}
}

// ---------------- CLIENT WRITE PUMP ----------------

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		// close send channel to signal readers
		c.hub.unregister <- c
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// hub closed the channel
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
