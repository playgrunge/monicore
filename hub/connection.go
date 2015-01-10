package hub

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/playgrunge/monicore/api"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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
	send chan []byte

	listCurrentAPI map[string]struct{}
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
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		var subscribedAPI []string
		json.Unmarshal(message, &subscribedAPI)

		for api := range subscribedAPI {
			c.listCurrentAPI[subscribedAPI[api]] = struct{}{}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
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
			var apiMessage ApiMessage
			json.Unmarshal(message, &apiMessage)

			_, ok2 := c.listCurrentAPI[apiMessage.Type]

			if ok2 {
				if err := c.write(websocket.TextMessage, message); err != nil {
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
	c := &connection{send: make(chan []byte, 256), ws: ws, listCurrentAPI: make(map[string]struct{})}
	h.register <- c
	go c.writePump()
	c.readPump()
}

func WsSend(w http.ResponseWriter, r *http.Request) {

	message := ApiMessage{}
	apiType := "chat"

	if r.FormValue("m") != "" {
		message = ApiMessage{apiType, r.FormValue("m")}
	} else {
		message = ApiMessage{apiType, "New message send from the server"}
	}

	messageJSON, _ := json.Marshal(message)
	h.broadcast <- messageJSON
}

func WsSendJSON(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://api.hockeystreams.com/Scores?key=" + api.GetConfig().Hockeystream.Key)
	if err != nil {
		log.Println(err)
		return
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	var hockeyData interface{}
	json.Unmarshal(robots, &hockeyData)

	apiMessage := ApiMessage{"hockey", hockeyData}
	messageToSend, _ := json.Marshal(apiMessage)

	h.broadcast <- messageToSend
}

type ApiMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
