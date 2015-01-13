package api

import (
	"io/ioutil"
	"log"
	"net/http"
)

type AirportApi struct{}

func (a *AirportApi) SendApi(w http.ResponseWriter, r *http.Request) {
	res, err := a.GetApi()
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (a *AirportApi) GetApi() ([]byte, error) {
	res, err := http.Get("http://api.flightstats.com/flex/airports/rest/v1/json/iata/YUL?appId=" + GetConfig().Flightstats.AppId + "&appKey=" + GetConfig().Flightstats.AppKey + "&utc=false&numHours=1&maxFlights=5")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return robots, err
}

func (a *AirportApi) updateData(data []byte) {

}
