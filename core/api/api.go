package api

import (
	"net/http"
)

type ApiRequest interface {
	SendApi(w http.ResponseWriter, r *http.Request)
	GetApi() ([]byte, error)
}
