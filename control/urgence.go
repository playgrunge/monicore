package control

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/playgrunge/monicore/core/api"
	"github.com/playgrunge/monicore/core/scrape"
	"log"
	"net/http"
	"strings"
)

type UrgenceApi struct {
	api.ApiRequest
	scrape.ScrapeRequest
}

func (u *UrgenceApi) Scrape(doc *goquery.Document) map[string]interface{} {
	var data = map[string]interface{}{}
	doc.Find("#oReportCell table tbody tr:nth-child(1) td div table tbody tr:nth-child(4) td:nth-child(2) table tbody tr").Each(func(i int, s *goquery.Selection) {
		urgence := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
		tauxOccupation := strings.TrimSpace(s.Find("td:nth-child(2)").Text())

		if i > 2 && urgence != "Total" && urgence != "Sous-total" && tauxOccupation != "" {
			data[urgence] = map[string]interface{}{
				"taux-occupation": tauxOccupation,
			}
		}
	})

	return data
}

const UrgenceName = "urgence"

func (u *UrgenceApi) SendApi(w http.ResponseWriter, r *http.Request) {
	res, err := u.GetApi()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(res)
}

func (u *UrgenceApi) GetApi() ([]byte, error) {
	doc, err := goquery.NewDocument("http://agence.santemontreal.qc.ca/fileadmin/asssm/rapports/urgence_quotidien_media.html")
	if err != nil {
		log.Fatal(err)
	}
	data := u.Scrape(doc)
	robots, err := json.Marshal(&data)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return robots, nil
}
