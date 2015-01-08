package api

import (
	"io/ioutil"
	"log"
	"net/http"
)

type HockeyApi struct{}

func (a *HockeyApi) SendApi(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://api.hockeystreams.com/Scores?key=" + GetConfig().Hockeystream.Key)
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(robots)
}
