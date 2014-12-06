package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/playgrunge/monicore/api"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/{key}", renderApi)
	r.HandleFunc("/pool", renderLongPool)
	//r.PathPrefix("/").Handler(NoCacheFileServer(http.Dir("./doc/")))
	r.PathPrefix("/").Handler(NoCacheFileServer(http.Dir("./app/")))
	http.Handle("/", r)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func renderApi(w http.ResponseWriter, r *http.Request) {
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

var lpchan = make(chan chan string)

func renderLongPool(w http.ResponseWriter, r *http.Request) {
	timeout, err := strconv.Atoi(r.URL.Query().Get("timeout"))
	if err != nil || timeout > 180000 || timeout < 0 {
		timeout = 60000 // default timeout is 60 seconds
	}
	var myRequestChan = make(chan string)

	select {
	case lpchan <- myRequestChan:
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		return
	}

	w.Write([]byte(<-myRequestChan))
}

// define global map;
var routes = map[string]interface{}{
	"hello":   hello_api,
	"bye":     bye_api,
	"test":    test_api,
	"json":    json_api,
	"hockey":  new(api.HockeyApi),
	"airport": new(api.AirportApi),
}

func hello_api(w http.ResponseWriter, r *http.Request) {
Loop:
	for {
		select {
		case clientchan := <-lpchan:
			clientchan <- "hello, client!"
			break
		default:
			break Loop
		}
	}
	w.Write([]byte("Hello World"))
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

type Message struct {
	Name string
	Body string
	Time int64
}
