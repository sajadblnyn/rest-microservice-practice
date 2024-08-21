package handlers

import (
	"errors"
	"net/http"

	"github.com/sajadblnyn/rest-microservice-practice/data"
)

type Product struct {
}

func (p *Product) ServeHTTP(rs http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rs, errors.New("method not allowed").Error(), http.StatusMethodNotAllowed)
		return
	}
	getProducts(rs, r)
}

func getProducts(rs http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()

	err := products.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
	}
}
