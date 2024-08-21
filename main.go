package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sajadblnyn/rest-microservice-practice/handlers"
)

func main() {

	m := http.NewServeMux()
	m.Handle("/", &handlers.Product{})
	s := &http.Server{Addr: ":8080",
		Handler: m, ReadTimeout: 1 * time.Second, WriteTimeout: 1 * time.Second, IdleTimeout: 100 * time.Second}
	go s.ListenAndServe()

	sigchan := make(chan os.Signal)

	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)

	sig := <-sigchan

	fmt.Println("shotting down server by signal: ", sig)

	c, _ := context.WithTimeout(context.Background(), time.Second*30)
	s.Shutdown(c)
}
