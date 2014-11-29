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
    
    if(key == "hello"){
        w.Write([]byte("Hello World"))
    }else if(key == "bye"){
        w.Write([]byte("Au revoir!!!"))
    }else if(key == "test"){
        w.Write([]byte("test"))
    }else if(key == "json"){
        m := Message{"Alice", "Hello", 1294706395881547000}
        b, _ := json.Marshal(m)
        w.Header().Set("Content-Type", "application/json")
        w.Write(b)
    }else{
        notFound(w,r)
    }
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
