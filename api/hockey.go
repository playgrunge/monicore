package api

import (
	"io/ioutil"
	"log"
	"net/http"
)

type HockeyApi struct{}

func (a *HockeyApi) SendApi(w http.ResponseWriter, r *http.Request) {
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
