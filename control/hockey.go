package control

import (
	"github.com/playgrunge/monicore/core/api"
	"io/ioutil"
	"log"
	"net/http"
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

	return robots, nil
}
