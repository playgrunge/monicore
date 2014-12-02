package api

import (
	"io/ioutil"
	"log"
	"net/http"
)

type AirportApi struct{}

func (a *AirportApi) SendApi(w http.ResponseWriter, r *http.Request) {
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
