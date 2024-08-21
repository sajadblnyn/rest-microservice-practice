package data

import (
	"encoding/json"
	"io"
)

type Product struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Price int64  `json:"price"`
}
type Products []*Product

func GetProducts() Products {
	products := Products{
		&Product{Id: 1, Title: "phone", Price: 200000},
		&Product{Id: 2, Title: "laptop", Price: 300000},
	}
	return products
}

func (p Products) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}
