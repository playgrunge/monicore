package api

import (
	"io/ioutil"
	"log"
	"net/http"
)

type AirportApi struct{}

func (a *AirportApi) SendApi(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://api.flightstats.com/flex/airports/rest/v1/json/iata/YUL?appId=" + getConfig().Flightstats.AppId + "&appKey=" + getConfig().Flightstats.AppKey + "&utc=false&numHours=1&maxFlights=5")
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
