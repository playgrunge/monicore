package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/playgrunge/monicore/api"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/{key}", renderApi)
	r.PathPrefix("/").Handler(NoCacheFileServer(http.Dir("./doc/")))
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

// define global map;
var routes = map[string]interface{}{
	"hello":   hello_api,
	"bye":     bye_api,
	"test":    test_api,
	"json":    json_api,
	"hockey":  new(api.HockeyApi),
	"airport": airport_api,
}

func hello_api(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
func bye_api(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Au revoir!!!"))
}
func test_api(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
func json_api(w http.ResponseWriter, r *http.Request) {
	//m := Message{"Alice", "Hello", 1294706395881547000}
	m := map[string]Message{
		"0": Message{"Alice", "Hello", 1294706395881547000},
		"1": Message{"Alice", "Hello", 1294706395881547000},
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func hockey_api(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://api.hockeystreams.com/Scores?key=f8788882ac0e9a9091c3985ce12fae82")
	if err != nil {
		log.Println(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(robots)
}

func airport_api(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("https://api.flightstats.com/flex/flightstatus/rest/v2/json/airport/status/YUL/dep/2014/11/30/10?appId=e454b3d5&appKey=6a2556db5129b9f57723eb368d34ae32&utc=false&numHours=1&maxFlights=5")
	if err != nil {
		log.Println(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(robots)
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
