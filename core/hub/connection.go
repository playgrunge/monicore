package hub

import (
	"github.com/gorilla/websocket"
	"github.com/playgrunge/monicore/helper"
	"log"
	"net/http"
	"sync"
	"time"
)

var mutexMessageTypes = &sync.Mutex{}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *Message

	messageTypes map[string]struct{}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))

	c.ws.SetPongHandler(
		func(string) error {
			c.ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

	for {
		clientMessageTypes := []string{}
		err := c.ws.ReadJSON(&clientMessageTypes)
		if err != nil {
			break
		}

		newMessageType := make(map[string]struct{})

		for t := range clientMessageTypes {
			newMessageType[clientMessageTypes[t]] = struct{}{}
		}

		newTypes := []string{}
		mutexMessageTypes.Lock()
		newTypes = helper.CompareMapKey(newMessageType, c.messageTypes)
		c.messageTypes = newMessageType
		mutexMessageTypes.Unlock()

		h.ReceiveNewTypes <- &PairConTypes{c, newTypes}
	}
}

func (c *connection) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}

func (c *connection) writeJSON(message interface{}) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(message)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:

			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			mutexMessageTypes.Lock()
			_, ok2 := c.messageTypes[message.Type]
			mutexMessageTypes.Unlock()

			if ok2 {
				if err := c.writeJSON(message); err != nil {
					return
				}
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serverWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan *Message, 256), ws: ws, messageTypes: make(map[string]struct{})}
	h.register <- c
	go c.writePump()
	c.readPump()
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
