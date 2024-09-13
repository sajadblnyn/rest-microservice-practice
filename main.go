package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"google.golang.org/grpc"

	"github.com/gorilla/mux"
	"github.com/sajadblnyn/rest-microservice-practice/data"
	"github.com/sajadblnyn/rest-microservice-practice/handlers"
	"github.com/sajadblnyn/rest-microservice-practice/middlewares"
	"github.com/sajadblnyn/rest-microservice-practice/protos/currency/protos/currency"
)

func main() {
	cc, err := grpc.Dial("127.0.0.1:9092", grpc.WithInsecure())
	if err != nil {
		log.Println("error in connecting ro currency service by grpc:", err)
	}
	gc := currency.NewCurrencyClient(cc)

	dp := data.NewProductDB(gc)
	productHandler := handlers.NewProductHandler(dp)
	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productHandler.GetProducts).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/", productHandler.GetProducts)
	getRouter.HandleFunc("/{id:[0-9]+}", productHandler.GetProductById).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/{id:[0-9]+}", productHandler.GetProductById)

	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", productHandler.UpdateProduct)
	putRouter.Use(middlewares.ProductMiddlewareFactory)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productHandler.AddProducts)
	postRouter.Use(middlewares.ProductMiddlewareFactory)
	// m.Handle("/", &handlers.Product{})

	uploaderRouter := router.Methods(http.MethodPost).Subrouter()
	uploaderRouter.HandleFunc("/storage", handlers.Upload)

	storageRouter := router.Methods(http.MethodGet).Subrouter()
	storageRouter.Handle("/storage/{filename:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\\.[a-z]{3}}", http.StripPrefix("/storage/", http.FileServer(http.Dir("./storage/"))))
	storageRouter.Use(middlewares.GzipMiddlewareFactory)

	co := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"http://127.0.0.1:8080", "http://localhost:8080"}))

	s := &http.Server{Addr: ":8080",
		Handler: co(router), ReadTimeout: 1 * time.Second, WriteTimeout: 1 * time.Second, IdleTimeout: 100 * time.Second}
	go s.ListenAndServe()

	sigchan := make(chan os.Signal)

	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)

	sig := <-sigchan

	fmt.Println("shotting down server by signal: ", sig)

	c, _ := context.WithTimeout(context.Background(), time.Second*30)
	s.Shutdown(c)
}
