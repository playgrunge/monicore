package control

import (
	"encoding/json"
	"github.com/playgrunge/monicore/core/api"
	"github.com/playgrunge/monicore/core/config"
	"github.com/playgrunge/monicore/core/hub"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

type HockeyApi struct {
	api.ApiRequest
}

const HockeyName = "hockey"

func (h *HockeyApi) SendApi(w http.ResponseWriter, r *http.Request) {
	res, err := h.GetApi()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (h *HockeyApi) GetApi() ([]byte, error) {
	res, err := http.Get("http://api.hockeystreams.com/Scores?key=" + config.GetConfig().Hockeystream.Key)
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

	h.updateData(robots)
	return robots, nil
}

func (h *HockeyApi) GetData() map[string]interface{} {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB("monicore").C(HockeyName)
	var r map[string]interface{}
	err = c.Find(nil).Sort("-timeStamp").Limit(1).One(&r)
	delete(r, "_id")
	delete(r, "timeStamp")

	return r
}

func (h *HockeyApi) updateData(data []byte) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB("monicore").C(HockeyName)
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
		message := hub.Message{HockeyName, d}
		hub.GetHub().Broadcast <- &message
		log.Println("Data updated")
	}

}
