package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/playgrunge/monicore/api"
	"github.com/playgrunge/monicore/hub"
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
	go run()

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func run() {
	ticker := time.NewTicker(time.Second * 10)
	log.Println("Ticker started")
	for _ = range ticker.C {
		if val, ok := routes[api.HockeyName]; ok {
			if t, ok := val.(api.ApiRequest); ok {
				t.GetApi()
			}
		}
	}
	defer ticker.Stop()
	log.Println("Ticker stopped")
}

func renderApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("RequestURI: " + r.Host + r.RequestURI)
	key := mux.Vars(r)["key"]

	if val, ok := routes[key]; ok {
		if t, ok := val.(api.ApiRequest); ok {
			t.SendApi(w, r)
		} else {
			val.(func(http.ResponseWriter, *http.Request))(w, r)
		}
	} else {
		notFound(w, r)
	}
}

// define global map;
var routes = map[string]interface{}{
	api.HockeyName:  new(api.HockeyApi),
	api.AirportName: new(api.AirportApi),
	api.WeatherName: new(api.WeatherApi),
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
	hockeyApi := new(api.HockeyApi)
	res, err := hockeyApi.GetApi()
	if err != nil {
		return
	}

	var hockeyData interface{}
	json.Unmarshal(res, &hockeyData)

	message := hub.Message{api.HockeyName, hockeyData}

	h.Broadcast <- &message
}

type Message struct {
	Name string
	Body string
	Time int64
}
