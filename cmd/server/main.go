package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nicholaslam/palindrome-service/internal/endpoint"
	"github.com/nicholaslam/palindrome-service/internal/service"
	"github.com/nicholaslam/palindrome-service/internal/store"
	"github.com/nicholaslam/palindrome-service/internal/transport"
)

func main() {
	httpAddr := flag.String("http-addr", ":8080", "HTTP listen address")
	strictPalindrome := flag.Bool("strict-palindrome", true, "Use strict definition of a palindrome")
	flag.Parse()

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

	// Duplicate routes to match trailing slash without redirecting.
	r.Methods("POST").Path("/messages").Handler(createHandler)
	r.Methods("POST").Path("/messages/").Handler(createHandler)
	r.Methods("GET").Path("/messages/{id}").Handler(readHandler)
	r.Methods("GET").Path("/messages/{id}/").Handler(readHandler)
	r.Methods("GET").Path("/messages").Handler(listHandler)
	r.Methods("GET").Path("/messages/").Handler(listHandler)
	r.Methods("DELETE").Path("/messages/{id}").Handler(deleteHandler)
	r.Methods("DELETE").Path("/messages/{id}/").Handler(deleteHandler)

	r.NotFoundHandler = http.HandlerFunc(notFound)

	log.Printf("http-addr %s", *httpAddr)
	log.Printf("strict-palindrome %t", *strictPalindrome)
	http.ListenAndServe(*httpAddr, r)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
