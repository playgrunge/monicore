package api

import (
	"encoding/json"
	"github.com/playgrunge/monicore/hub"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

type WeatherApi struct{}

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

	a.updateData(robots)
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

	a.updateData(robots)
	return robots, nil
}

func (a *WeatherApi) GetData() map[string]interface{} {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB("monicore").C("weather")
	var r map[string]interface{}
	err = c.Find(nil).Sort("-timeStamp").Limit(1).One(&r)
	delete(r, "_id")
	delete(r, "timeStamp")

	return r
}

func (a *WeatherApi) updateData(data []byte) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB("monicore").C("weather")
	var r map[string]interface{}
	err = c.Find(nil).Sort("-timeStamp").Limit(1).One(&r)
	delete(r, "_id")
	delete(r, "timeStamp")

	var d map[string]interface{}
	json.Unmarshal(data, &d)

	eq := reflect.DeepEqual(r, d)
	if !eq {
		d["timeStamp"] = time.Now()
		err = c.Insert(d)
		message := hub.Message{"weather", d}
		hub.GetHub().Broadcast <- &message
		log.Println("Data updated")
	}

}
