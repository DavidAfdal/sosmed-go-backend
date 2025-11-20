package socket

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	UserID string
	Send   chan []byte
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Connection]bool
	rdb     *redis.Client
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		clients: make(map[string]map[*Connection]bool),
		rdb:     rdb,
	}
}

func (h *Hub) Register(userID string, conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[userID]; ok {
		h.clients[userID] = make(map[*Connection]bool)
	}
	h.clients[userID][conn] = true
}

func (h *Hub) Unregister(userID string, conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[userID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.clients, userID)
		}
	}
}

func (h *Hub) SendToUser(userID string, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if conns, ok := h.clients[userID]; ok {
		for c := range conns {
			c.Send(msg)
		}
	}
}
