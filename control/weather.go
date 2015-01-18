package control

import (
	"github.com/playgrunge/monicore/core/api"
	"io/ioutil"
	"log"
	"net/http"
)

type WeatherApi struct {
	api.ApiRequest
}

const WeatherName = "weather"

func (a *WeatherApi) SendApi(w http.ResponseWriter, r *http.Request) {
	typ := r.FormValue("t")
	if len(typ) == 0 {
		typ = "default"
	}
	var res []byte
	var err error
	switch typ {
	case "forecast":
		res, err = a.GetForecast()
	default:
		res, err = a.GetApi()
	}

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (a *WeatherApi) GetApi() ([]byte, error) {
	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=montreal&APPID=" + GetConfig().Openweathermap.AppId)
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

	return robots, nil
}

func (a *WeatherApi) GetForecast() ([]byte, error) {
	res, err := http.Get("http://api.openweathermap.org/data/2.5/forecast/city?q=montreal&APPID=" + GetConfig().Openweathermap.AppId)
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

	return robots, nil
}
