package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"github.com/Igor-K/mgo"
	"reflect"
)

type HockeyApi struct{}

func (h *HockeyApi) SendApi(w http.ResponseWriter, r *http.Request) {
	res, err := h.GetApi()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (h *HockeyApi) GetApi() ([]byte, error) {
	res, err := http.Get("http://api.hockeystreams.com/Scores?key=" + GetConfig().Hockeystream.Key)
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

func (h *HockeyApi) GetData() (map[string]interface{}) {
	session, err := mgo.Dial("localhost")
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    c := session.DB("monicore").C("hockey")
    var r map[string]interface{}
    err = c.Find(nil).Sort("$natural:-1").Limit(1).One(&r);
    delete(r, "_id")
    return r
}

func (h *HockeyApi) updateData(data []byte) {
	session, err := mgo.Dial("localhost")
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    c := session.DB("monicore").C("hockey")
    var r map[string]interface{}
    err = c.Find(nil).Sort("$natural:-1").Limit(1).One(&r);
    delete(r, "_id")

	log.Println("Data compared...")
	var d map[string]interface{}
	json.Unmarshal(data, &d)
	eq := reflect.DeepEqual(r, d)
	if eq {
		log.Println("Data identical")
	} else {
		err = c.Insert(d)
		log.Println("Data updated")
	}
	
}