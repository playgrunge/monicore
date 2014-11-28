package main

import (
	"log"
	"net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/api/{key}", renderApi)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./doc/")))
    http.Handle("/", r)

    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)
}

func renderApi(w http.ResponseWriter, r *http.Request){
    log.Println("RequestURI: " + r.Host+r.RequestURI)
    key := mux.Vars(r)["key"]
    
    if(key == "hello"){
        w.Write([]byte("Hello World"))
    }else{
        notFound(w,r)
    }
}

func notFound(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(400)
    w.Write([]byte("400: Bad Request"))
}
