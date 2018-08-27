package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nicholaslam/palindrome-service/internal/endpoint"
	"github.com/nicholaslam/palindrome-service/internal/service"
	"github.com/nicholaslam/palindrome-service/internal/store"
	"github.com/nicholaslam/palindrome-service/internal/transport"
)

func main() {
	defaultHTTPAddr := ":8080"
	defaultStrictPalindrome := true

	httpAddr := flag.String("http-addr", defaultHTTPAddr, "HTTP listen address")
	strictPalindrome := flag.Bool("strict-palindrome", defaultStrictPalindrome, "Use strict definition of a palindrome")
	flag.Parse()

	envHTTPAddr := os.Getenv("HTTP_ADDR")
	if *httpAddr == defaultHTTPAddr && envHTTPAddr != "" {
		httpAddr = &envHTTPAddr
	}

	var err error
	envStrictPalindrome := os.Getenv("STRICT_PALINDROME")
	if *strictPalindrome == defaultStrictPalindrome && envStrictPalindrome != "" {
		*strictPalindrome, err = strconv.ParseBool(envStrictPalindrome)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	store := store.NewTempStore()
	service := service.NewService(store, *strictPalindrome)

	createEndpoint := endpoint.MakeCreateEndpoint(service)
	readEndpoint := endpoint.MakeReadEndpoint(service)
	listEndpoint := endpoint.MakeListEndpoint(service)
	deleteEndpoint := endpoint.MakeDeleteEndpoint(service)

	createHandler := transport.MakeCreateHTTPHandler(createEndpoint)
	readHandler := transport.MakeReadHTTPHandler(readEndpoint)
	listHandler := transport.MakeListHTTPHandler(listEndpoint)
	deleteHandler := transport.MakeDeleteHTTPHandler(deleteEndpoint)

	r := mux.NewRouter()
	r = r.PathPrefix("/api/v1/").Subrouter()

	// Duplicate route definitions to match trailing slash without redirecting.
	r.Methods("POST").Path("/messages").Handler(createHandler)
	r.Methods("POST").Path("/messages/").Handler(createHandler)
	r.Methods("GET").Path("/messages/{id}").Handler(readHandler)
	r.Methods("GET").Path("/messages/{id}/").Handler(readHandler)
	r.Methods("GET").Path("/messages").Handler(listHandler)
	r.Methods("GET").Path("/messages/").Handler(listHandler)
	r.Methods("DELETE").Path("/messages/{id}").Handler(deleteHandler)
	r.Methods("DELETE").Path("/messages/{id}/").Handler(deleteHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)

	srv := http.Server{
		Addr:    *httpAddr,
		Handler: r,
	}

	done := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
		close(done)
	}()
	log.Println("server started", "http-addr", *httpAddr, "strict-palindrome", *strictPalindrome)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
	<-done
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
