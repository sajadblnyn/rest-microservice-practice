package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sajadblnyn/rest-microservice-practice/data"
	"github.com/sajadblnyn/rest-microservice-practice/responses"
)

func ProductMiddlewareFactory(next http.Handler) http.Handler {
	return &ProductMiddleware{nextHandler: next}
}

type ProductMiddleware struct{ nextHandler http.Handler }

type ProductKeyType string

var ProductKey ProductKeyType = "ProductKey"

func (p *ProductMiddleware) ServeHTTP(rs http.ResponseWriter, rq *http.Request) {
	prod := &data.Product{}

	err := prod.FromJson(rq.Body)
	if err != nil {
		http.Error(rs, err.Error(), http.StatusBadRequest)
		return
	}

	err = prod.Validate()
	response := responses.MakeResponse(err, nil)

	if err != nil {
		jsonRes, err := json.Marshal(response)
		if err != nil {
			http.Error(rs, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(rs, string(jsonRes), http.StatusBadRequest)
		return
	}
	ctx := context.WithValue(rq.Context(), ProductKey, prod)
	rq = rq.WithContext(ctx)

	p.nextHandler.ServeHTTP(rs, rq)

}
