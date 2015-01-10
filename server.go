package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/playgrunge/monicore/api"
	"github.com/playgrunge/monicore/hub"
	"io/ioutil"
	"log"
	"net/http"
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

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
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
	"bye":     bye_api,
	"test":    test_api,
	"json":    json_api,
	"hockey":  new(api.HockeyApi),
	"airport": new(api.AirportApi),
}

func bye_api(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Au revoir!!!"))
}
func test_api(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
func json_api(w http.ResponseWriter, r *http.Request) {
	m := map[string]Message{
		"0": Message{"Alice", "Hello", 1294706395881547000},
		"1": Message{"Bob", "Bye", 1294706595681746000},
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
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

	message := hub.Message{"hockey", hockeyData}

	h.Broadcast <- &message
}

type Message struct {
	Name string
	Body string
	Time int64
}
