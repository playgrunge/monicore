package hub

import (
	"log"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]struct{}

	// Inbound messages from the connections.
	Broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	Broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]struct{}),
}

func GetHub() *hub {
	return &h
}
func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			log.Println("Register...")
			h.connections[c] = struct{}{}
		case c := <-h.unregister:
			log.Println("Unregister...")
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.Broadcast:
			log.Println("Send data...")
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
