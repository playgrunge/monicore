package service

import (
	"encoding/json"
	"github.com/playgrunge/monicore/core/hub"
	"gopkg.in/mgo.v2"
	"log"
	"reflect"
	"time"
)

func UpdateNewData(dataType string, data []byte) bool {
	isDataUpdated := false
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB("monicore").C(dataType)
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
		message := hub.Message{dataType, d}
		hub.GetHub().Broadcast <- &message
		isDataUpdated = true
		log.Println("Data updated")
	} else {
		log.Println("Not updated")
	}
	return isDataUpdated
}

func GetLastData(dataType string) map[string]interface{} {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB("monicore").C(dataType)
	var r map[string]interface{}
	err = c.Find(nil).Sort("-timeStamp").Limit(1).One(&r)
	delete(r, "_id")
	delete(r, "timeStamp")

	return r
}
