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
	Broadcast chan *Message

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection

	ReceiveNewTypes chan *PairConTypes

	SendToConnection chan *PairConMessage
}

var h = hub{
	Broadcast:        make(chan *Message),
	register:         make(chan *connection),
	unregister:       make(chan *connection),
	connections:      make(map[*connection]struct{}),
	ReceiveNewTypes:  make(chan *PairConTypes),
	SendToConnection: make(chan *PairConMessage),
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
			log.Println("Broadcast data...")
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		case s := <-h.SendToConnection:
			c := s.Con
			log.Println("Send message to connection")
			if _, ok := h.connections[c]; ok {
				select {
				case c.send <- s.Message:
				default:
					close(c.send)
					delete(h.connections, c)
				}

			}
		}
	}
}

type PairConTypes struct {
	Con   *connection
	Types []string
}

type PairConMessage struct {
	Con     *connection
	Message *Message
}
