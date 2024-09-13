package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sajadblnyn/rest-microservice-practice/data"
	"github.com/sajadblnyn/rest-microservice-practice/middlewares"
)

type Product struct {
	dp *data.ProductDB
}

// func (p *Product) ServeHTTP(rs http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		getProducts(rs, r)
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		addProducts(rs, r)
// 		return
// 	}

// 	if r.Method == http.MethodPut {
// 		pathId := strings.TrimPrefix(r.URL.Path, "/")
// 		id, err := strconv.Atoi(pathId)
// 		if err != nil {
// 			http.Error(rs, "path parameter id must be integer", http.StatusBadRequest)
// 			return
// 		}

// 		updateProduct(id, rs, r)
// 		return
// 	}

// 	http.Error(rs, errors.New("method not allowed").Error(), http.StatusMethodNotAllowed)

// }

func NewProductHandler(dp *data.ProductDB) *Product {
	return &Product{dp: dp}
}

func (p *Product) GetProducts(rs http.ResponseWriter, r *http.Request) {

	cr := r.URL.Query().Get("currency")
	if cr == "" {
		cr = "USD"
	}
	products, err := p.dp.GetProducts(cr)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}
	err = products.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *Product) GetProductById(rs http.ResponseWriter, r *http.Request) {

	cr := r.URL.Query().Get("currency")
	if cr == "" {
		cr = "USD"
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rs, err.Error(), http.StatusBadRequest)
		return
	}
	product, err := p.dp.GetProductById(id, cr)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}
	err = product.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *Product) UpdateProduct(rs http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rs, err.Error(), http.StatusBadRequest)
		return
	}

	prod := (r.Context().Value(middlewares.ProductKey).(*data.Product))

	err = p.dp.UpdateProduct(id, prod)

	if err == data.NotfoundError {
		http.Error(rs, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}
	err = prod.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (p *Product) AddProducts(rs http.ResponseWriter, r *http.Request) {

	prod := (r.Context().Value(middlewares.ProductKey).(*data.Product))

	p.dp.AddProduct(prod)

	err := prod.ToJson(rs)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusInternalServerError)
		return
	}

}
