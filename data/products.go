package data

import (
	"encoding/json"
	"errors"
	"io"
)

type Product struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Price int64  `json:"price"`
}
type Products []*Product

var products = &Products{
	&Product{Id: 1, Title: "phone", Price: 200000},
	&Product{Id: 2, Title: "laptop", Price: 300000},
}

var NotfoundError error = errors.New("product not found")

func GetProducts() *Products {
	return products
}

func AddProduct(p *Product) {
	prods := *products
	p.Id = prods[len(prods)-1].Id + 1

	prods = append(prods, p)
	products = &prods
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.Id = id
	(*products)[pos] = p

	return nil
}

func findProduct(id int) (p *Product, pos int, err error) {
	pos = -1
	for i, v := range *products {
		if v.Id == id {
			pos = i
			p = v
		}
	}
	if pos == -1 {
		return nil, pos, NotfoundError
	}
	return p, pos, nil

}

func (p *Products) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJson(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *Product) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}
