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

var (
	defaultHTTPAddr         = ":8080"
	defaultStrictPalindrome = true
)

type config struct {
	httpAddr         string
	strictPalindrome bool
}

func main() {
	cfg := parseConfig(os.Args)

	store := store.NewTempStore()
	service := service.NewService(store, cfg.strictPalindrome)

	createEndpoint := endpoint.MakeCreateEndpoint(service)
	readEndpoint := endpoint.MakeReadEndpoint(service)
	listEndpoint := endpoint.MakeListEndpoint(service)
	deleteEndpoint := endpoint.MakeDeleteEndpoint(service)

	createHandler := transport.MakeCreateHTTPHandler(createEndpoint)
	readHandler := transport.MakeReadHTTPHandler(readEndpoint)
	listHandler := transport.MakeListHTTPHandler(listEndpoint)
	deleteHandler := transport.MakeDeleteHTTPHandler(deleteEndpoint)

	// Duplicate route definitions to match trailing slash without redirecting.
	r := mux.NewRouter()
	r.Methods("GET").Path("/healthz").HandlerFunc(healthz)
	r.Methods("GET").Path("/healthz/").HandlerFunc(healthz)

	s := r.PathPrefix("/api/v1/").Subrouter()
	s.Methods("POST").Path("/messages").Handler(createHandler)
	s.Methods("POST").Path("/messages/").Handler(createHandler)
	s.Methods("GET").Path("/messages/{id}").Handler(readHandler)
	s.Methods("GET").Path("/messages/{id}/").Handler(readHandler)
	s.Methods("GET").Path("/messages").Handler(listHandler)
	s.Methods("GET").Path("/messages/").Handler(listHandler)
	s.Methods("DELETE").Path("/messages/{id}").Handler(deleteHandler)
	s.Methods("DELETE").Path("/messages/{id}/").Handler(deleteHandler)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	srv := http.Server{
		Addr:    cfg.httpAddr,
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

	log.Println("http-addr", cfg.httpAddr, "strict-palindrome", cfg.strictPalindrome)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
		if err != http.ErrServerClosed {
			close(done)
		}
	}

	<-done
}

func parseConfig(args []string) config {
	fsName := args[0]
	fsArgs := args[1:]

	fs := flag.NewFlagSet(fsName, flag.ExitOnError)
	httpAddr := fs.String("http-addr", defaultHTTPAddr, "HTTP listen address")
	strictPalindrome := fs.Bool("strict-palindrome", defaultStrictPalindrome, "Use strict definition of a palindrome")
	fs.Parse(fsArgs)

	envHTTPAddr := os.Getenv("HTTP_ADDR")
	if *httpAddr == defaultHTTPAddr && envHTTPAddr != "" {
		*httpAddr = envHTTPAddr
	}

	var err error
	envStrictPalindrome := os.Getenv("STRICT_PALINDROME")
	if *strictPalindrome == defaultStrictPalindrome && envStrictPalindrome != "" {
		*strictPalindrome, err = strconv.ParseBool(envStrictPalindrome)
		if err != nil {
			fmt.Printf("invalid boolean value \"%s\" for STRICT_PALINDROME: %s\n", envStrictPalindrome, err.Error())
			os.Exit(2)
		}
	}

	return config{
		*httpAddr,
		*strictPalindrome,
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
