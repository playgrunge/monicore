package control

import (
	"github.com/playgrunge/monicore/core/api"
	"github.com/playgrunge/monicore/core/config"
	"io/ioutil"
	"log"
	"net/http"
)

type AirportApi struct {
	api.ApiRequest
}

const AirportName = "airport"

func (a *AirportApi) SendApi(w http.ResponseWriter, r *http.Request) {
	typ := r.FormValue("t")
	if len(typ) == 0 {
		typ = "default"
	}
	var res []byte
	var err error
	switch typ {
	case "delay":
		res, err = a.GetDelay()
	case "weather":
		res, err = a.GetWeather()
	default:
		res, err = a.GetApi()
	}

	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (a *AirportApi) GetApi() ([]byte, error) {
	res, err := http.Get("http://api.flightstats.com/flex/airports/rest/v1/json/iata/YUL?appId=" + config.GetConfig().Flightstats.AppId + "&appKey=" + config.GetConfig().Flightstats.AppKey + "&utc=false&numHours=1&maxFlights=5")
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

func (a *AirportApi) GetWeather() ([]byte, error) {
	res, err := http.Get("http://api.flightstats.com/flex/weather/rest/v1/json/all/YUL?appId=" + config.GetConfig().Flightstats.AppId + "&appKey=" + config.GetConfig().Flightstats.AppKey + "&utc=false&numHours=1&maxFlights=5")
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

func (a *AirportApi) GetDelay() ([]byte, error) {
	res, err := http.Get("http://api.flightstats.com/flex/delayindex/rest/v1/json/airports/YUL?appId=" + config.GetConfig().Flightstats.AppId + "&appKey=" + config.GetConfig().Flightstats.AppKey + "&utc=false&numHours=1&maxFlights=5")
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
