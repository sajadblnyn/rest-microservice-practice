package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

	if r.Method == http.MethodPut {
		pathId := strings.TrimPrefix(r.URL.Path, "/")
		id, err := strconv.Atoi(pathId)
		if err != nil {
			http.Error(rs, "path parameter id must be integer", http.StatusBadRequest)
			return
		}

		updateProduct(id, rs, r)
		return
	}

	http.Error(rs, errors.New("method not allowed").Error(), http.StatusMethodNotAllowed)

}

func getProducts(rs http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()

	err := products.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateProduct(id int, rs http.ResponseWriter, r *http.Request) {

	prod := &data.Product{}

	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusBadRequest)
		return
	}
	err = data.UpdateProduct(id, prod)

	if err == data.NotfoundError {
		http.Error(rs, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}
	prod.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}

}

func addProducts(rs http.ResponseWriter, r *http.Request) {

	prod := &data.Product{}

	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusBadRequest)
		return
	}
	data.AddProduct(prod)

	prod.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}

}
