package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sajadblnyn/rest-microservice-practice/handlers"
	"github.com/sajadblnyn/rest-microservice-practice/middlewares"
)

func main() {

	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", handlers.GetProducts)

	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", handlers.UpdateProduct)
	putRouter.Use(middlewares.ProductMiddlewareFactory)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", handlers.AddProducts)
	postRouter.Use(middlewares.ProductMiddlewareFactory)
	// m.Handle("/", &handlers.Product{})

	s := &http.Server{Addr: ":8080",
		Handler: router, ReadTimeout: 1 * time.Second, WriteTimeout: 1 * time.Second, IdleTimeout: 100 * time.Second}
	go s.ListenAndServe()

	sigchan := make(chan os.Signal)

	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)

	sig := <-sigchan

	fmt.Println("shotting down server by signal: ", sig)

	c, _ := context.WithTimeout(context.Background(), time.Second*30)
	s.Shutdown(c)
}
