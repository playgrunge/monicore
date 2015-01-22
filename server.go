package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/playgrunge/monicore/control"
	"github.com/playgrunge/monicore/core/api"
	"github.com/playgrunge/monicore/core/hub"
	"github.com/playgrunge/monicore/service"
	"log"
	"net/http"
	"time"
)

var h = hub.GetHub()

func main() {
	r := mux.NewRouter()
	http.HandleFunc("/websocket", hub.ServeWs)
	http.HandleFunc("/wsSend", wsSend)
	http.HandleFunc("/wsSendJSON", wsSendJSON)
	r.HandleFunc("/api/{key}", renderApi)
	r.PathPrefix("/").Handler(NoCacheFileServer(http.Dir("./doc/")))
	http.Handle("/", r)

	go h.Run()
	go runTaskUpdateData(control.HockeyName, time.Minute*10)
	go listenForNewTypes()

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func runTaskUpdateData(dataType string, tickerTime time.Duration) {
	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()
	log.Println("Ticker started")
	for _ = range ticker.C {
		if val, ok := routes[dataType]; ok {
			if t, ok := val.(api.ApiRequest); ok {
				broadcastMessageIfNew(dataType, t)
			}
		}
	}
	log.Println("Ticker stopped")
}

func listenForNewTypes() {
	for {
		c := <-h.ReceiveNewTypes
		go func() {
			for i := range c.Types {
				if val, ok := routes[c.Types[i]]; ok {
					if t, ok := val.(api.ApiRequest); ok {
						if val, _ := t.GetApi(); val != nil {
							var d map[string]interface{}
							json.Unmarshal(val, &d)
							message := hub.Message{c.Types[i], d}
							pairConMessage := &hub.PairConMessage{c.Con, &message}
							h.SendToConnection <- pairConMessage
						}
					}
				}
			}
		}()
	}
}

func broadcastMessageIfNew(dataType string, a api.ApiRequest) {
	if val, _ := a.GetApi(); val != nil {
		if isNew := service.UpdateNewData(dataType, val); isNew {
			lastData := service.GetLastData(dataType)
			message := hub.Message{dataType, lastData}
			hub.GetHub().Broadcast <- &message
		}
	}
}

func renderApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("RequestURI: " + r.Host + r.RequestURI)
	key := mux.Vars(r)["key"]

	if val, ok := routes[key]; ok {
		switch v := val.(type) {
		case api.ApiRequest:
			v.SendApi(w, r)
		default:
			v.(func(http.ResponseWriter, *http.Request))(w, r)
		}

	} else {
		notFound(w, r)
	}
}

// define global map;
var routes = map[string]interface{}{
	control.HockeyName:  new(control.HockeyApi),
	control.AirportName: new(control.AirportApi),
	control.WeatherName: new(control.WeatherApi),
	control.HydroName:   new(control.HydroApi),
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	w.Write([]byte("400: Bad Request"))
}

type noCacheFileServer struct {
	root http.FileSystem
}

func NoCacheFileServer(root http.FileSystem) http.Handler {
	return &noCacheFileServer{root}
}
func (n *noCacheFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.FileServer(n.root).ServeHTTP(w, r)
}

func wsSend(w http.ResponseWriter, r *http.Request) {

	message := hub.Message{}
	apiType := "chat"

	if r.FormValue("m") != "" {
		message = hub.Message{apiType, r.FormValue("m")}
	} else {
		message = hub.Message{apiType, "New message send from the server"}
	}

	h.Broadcast <- &message
}

func wsSendJSON(w http.ResponseWriter, r *http.Request) {
	if val, ok := routes[control.HockeyName]; ok {
		if t, ok := val.(api.ApiRequest); ok {
			broadcastMessageIfNew(control.HockeyName, t)
		}
	}
}

type Message struct {
	Name string
	Body string
	Time int64
}
