package main

import (
	"log"
	"net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/api/{key}", renderApi)
    r.PathPrefix("/").Handler(NoCacheFileServer(http.Dir("./doc/")))
    http.Handle("/", r)

    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)
}

func renderApi(w http.ResponseWriter, r *http.Request){
    log.Println("RequestURI: " + r.Host+r.RequestURI)
    key := mux.Vars(r)["key"]
    
    if val, ok := routes[key]; ok {
        val.(func(http.ResponseWriter, *http.Request))(w,r)
    }else{
        notFound(w,r)
    }
}
// define global map; initialize as empty with the trailing {}
var routes = map[string]interface{}{
    "hello": hello_api,
    "bye": bye_api,
    "test": test_api,
    "json": json_api,
}
func hello_api(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("Hello World"))
}
func bye_api(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("Au revoir!!!"))
}
func test_api(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("test"))
}
func json_api(w http.ResponseWriter, r *http.Request){
    m := Message{"Alice", "Hello", 1294706395881547000}
    b, _ := json.Marshal(m)
    w.Header().Set("Content-Type", "application/json")
    w.Write(b)
}

func notFound(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(400)
    w.Write([]byte("400: Bad Request"))
}

type noCacheFileServer struct{
    root http.FileSystem
}
func NoCacheFileServer(root http.FileSystem) http.Handler {
    return &noCacheFileServer{root}
}
func (n *noCacheFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
    http.FileServer(n.root).ServeHTTP(w,r)
}

type Message struct {
    Name string
    Body string
    Time int64
}
