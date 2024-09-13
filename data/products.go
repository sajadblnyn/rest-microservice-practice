package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/sajadblnyn/rest-microservice-practice/protos/currency/protos/currency"
)

type ProductDB struct {
	c            currency.CurrencyClient
	rates        map[string]float64
	subscription currency.Currency_SubscribeRatesClient
}

func NewProductDB(c currency.CurrencyClient) *ProductDB {
	pd := &ProductDB{c: c, rates: map[string]float64{}, subscription: nil}
	go pd.handleRatesUpdates()
	return pd
}

func (pd *ProductDB) handleRatesUpdates() {
	sub, err := pd.c.SubscribeRates(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	pd.subscription = sub
	for {
		rr, err := sub.Recv()
		fmt.Println("Recieved updated rate from server", "dest", rr.GetDestination().String(), "rate", rr.GetRate())

		if err != nil {
			fmt.Println(err)
		}
		pd.rates[rr.Destination.String()] = rr.GetRate()
	}
}

type Product struct {
	Id    int     `json:"id"`
	Title string  `json:"title" validate:"required,titleLength"`
	Price float64 `json:"price" validate:"gte=0,lte=2000000"`
}
type Products []*Product

var products = &Products{
	&Product{Id: 1, Title: "phone", Price: 200000},
	&Product{Id: 2, Title: "laptop", Price: 300000},
}

var NotfoundError error = errors.New("product not found")

func (pd *ProductDB) GetProducts(destCurrency string) (*Products, error) {

	rate, err := pd.getCurrencyRate(destCurrency)

	if err != nil {
		return nil, err
	}

	var prods Products
	for _, v := range *products {
		var prod Product = *v
		prod.Price = prod.Price * rate

		prods = append(prods, &prod)
	}
	return &prods, nil
}

func (pd *ProductDB) GetProductById(id int, destCurrency string) (*Product, error) {

	rate, err := pd.getCurrencyRate(destCurrency)
	if err != nil {
		return nil, err
	}

	p, _, err := pd.findProduct(id)
	if err != nil {
		return nil, err
	}

	np := *p
	np.Price = rate * np.Price
	return &np, nil
}

func (pd *ProductDB) AddProduct(p *Product) {
	prods := *products
	p.Id = prods[len(prods)-1].Id + 1

	prods = append(prods, p)
	products = &prods
}

func (pd *ProductDB) UpdateProduct(id int, p *Product) error {
	_, pos, err := pd.findProduct(id)
	if err != nil {
		return err
	}
	p.Id = id
	(*products)[pos] = p

	return nil
}

func (pd *ProductDB) getCurrencyRate(destCurrency string) (float64, error) {
	rate, ok := pd.rates[destCurrency]

	if ok {
		return rate, nil
	}

	rr, err := pd.c.GetRate(context.Background(), &currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_value[currency.Currencies_USD.String()]),
		Destination: currency.Currencies(currency.Currencies_value[destCurrency])})

	if err != nil {
		return 1, err
	}

	pd.rates[rr.Destination.String()] = rr.GetRate()

	err = pd.subscription.Send(&currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_value[currency.Currencies_USD.String()]),
		Destination: currency.Currencies(currency.Currencies_value[destCurrency])})

	return rr.Rate, err
}

func (pd *ProductDB) findProduct(id int) (p *Product, pos int, err error) {
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

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("titleLength", titleLengthValidation)

	return validate.Struct(p)
}

func titleLengthValidation(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) < 5 || len(fl.Field().String()) > 50 {
		return false
	}
	return true
}
