package control

import (
	"encoding/json"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"github.com/playgrunge/monicore/core/api"
	"log"
	"net/http"
)

type TwitterApi struct {
	api.ApiRequest
}

const TwitterName = "twitter"

func (t *TwitterApi) SendApi(w http.ResponseWriter, r *http.Request) {
	res, err := t.GetApi()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (t *TwitterApi) GetApi() ([]byte, error) {

	config := &oauth1a.ClientConfig{
		ConsumerKey:    GetConfig().Twitter.ConsumerKey,
		ConsumerSecret: GetConfig().Twitter.ConsumerSecret,
	}
	client := twittergo.NewClient(config, nil)

	req, err := http.NewRequest("GET", "/1.1/trends/place.json?id=3534", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := client.SendRequest(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	robots := []byte(res.ReadBody())

	var d []interface{}
	json.Unmarshal(robots, &d)

	result, _ := d[0].(map[string]interface{})

	return json.Marshal(result)
}
