package handlers

import (
	"errors"
	"net/http"

	"github.com/sajadblnyn/rest-microservice-practice/data"
)

type Product struct {
}

func (p *Product) ServeHTTP(rs http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getProducts(rs, r)
		return
	}

	if r.Method == http.MethodPost {
		addProducts(rs, r)
		return
	}

	http.Error(rs, errors.New("method not allowed").Error(), http.StatusMethodNotAllowed)

}

func getProducts(rs http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()

	err := products.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
	}
}

func addProducts(rs http.ResponseWriter, r *http.Request) {

	prod := &data.Product{}

	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusBadRequest)
	}
	data.AddProduct(prod)

	prod.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
	}

}
