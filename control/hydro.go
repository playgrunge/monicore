package control

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/playgrunge/monicore/core/api"
	"github.com/playgrunge/monicore/core/scrape"
	"log"
	"net/http"
	"regexp"
)

type HydroApi struct {
	api.ApiRequest
	scrape.ScrapeRequest
}

func (h *HydroApi) Scrape(doc *goquery.Document) map[string]interface{} {
	var intrClientsREGEX = regexp.MustCompile(`^[0-9 ]+`)
	var totalClientsREGEX = regexp.MustCompile(`[0-9 ]+$`)

	var data = map[string]interface{}{}
	doc.Find("div.service-on table tbody tr").Each(func(i int, s *goquery.Selection) {
		region := s.Find("td[scope=row] a").Text()
		interruptions := s.Find("td:nth-child(2)").Text()
		clients := s.Find("td:nth-child(3)").Text()
		intrClients := intrClientsREGEX.FindString(clients)
		totalClients := totalClientsREGEX.FindString(clients)

		data[region] = map[string]interface{}{
			"interruptions":      interruptions,
			"clientsInterrupted": intrClients,
			"totalClients":       totalClients,
		}
	})

	return data
}

const HydroName = "hydro"

func (h *HydroApi) SendApi(w http.ResponseWriter, r *http.Request) {
	res, err := h.GetApi()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(res)
}

func (h *HydroApi) GetApi() ([]byte, error) {
	doc, err := goquery.NewDocument("http://pannes.hydroquebec.com/pannes/bilan-interruptions-service")
	if err != nil {
		log.Fatal(err)
	}
	data := h.Scrape(doc)
	robots, err := json.Marshal(&data)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return robots, nil
}
