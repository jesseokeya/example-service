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
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nicholaslam/example-service/internal/endpoint"
	"github.com/nicholaslam/example-service/internal/service"
	"github.com/nicholaslam/example-service/internal/store"
	"github.com/nicholaslam/example-service/internal/transport"
)

const (
	dbName         = "palindromedb"
	collectionName = "messages"
)

var (
	defaultHTTPAddr         = ":8080"
	defaultStrictPalindrome = true
	defaultMongoURI         = ""
)

type config struct {
	httpAddr         string
	strictPalindrome bool
	mongoURI         string
}

func main() {
	cfg, err := parseConfig(os.Args)
	if err != nil {
		log.Println("error parsing config:", err)
		return
	}

	str := store.NewTempStore()
	if cfg.mongoURI != "" {
		client, err := mongo.NewClient(cfg.mongoURI)
		if err != nil {
			log.Println("error creating mongo client:", err)
			return
		}
		err = client.Connect(context.Background())
		if err != nil {
			log.Println("error connecting to mongo client:", err)
			return
		}
		str = store.NewMongoStore(client.Database(dbName).Collection(collectionName))
	}

	service := service.NewService(str, cfg.strictPalindrome)

	createEndpoint := endpoint.MakeCreateEndpoint(service)
	readEndpoint := endpoint.MakeReadEndpoint(service)
	listEndpoint := endpoint.MakeListEndpoint(service)
	deleteEndpoint := endpoint.MakeDeleteEndpoint(service)

	createHandler := transport.MakeCreateHTTPHandler(createEndpoint)
	readHandler := transport.MakeReadHTTPHandler(readEndpoint)
	listHandler := transport.MakeListHTTPHandler(listEndpoint)
	deleteHandler := transport.MakeDeleteHTTPHandler(deleteEndpoint)

	// Duplicate the route definitions to match trailing slash without redirecting.
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

func parseConfig(args []string) (config, error) {
	fsName := args[0]
	fsArgs := args[1:]

	fs := flag.NewFlagSet(fsName, flag.ExitOnError)
	httpAddr := fs.String("http-addr", defaultHTTPAddr, "HTTP listen address")
	strictPalindrome := fs.Bool("strict-palindrome", defaultStrictPalindrome, "Use strict definition of a palindrome")
	mongoURI := fs.String("mongo-uri", defaultMongoURI, "MongoDB connection string. Pass empty string to use in-memory database")
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
			err = fmt.Errorf(`invalid boolean value "%s" for STRICT_PALINDROME: %s`, envStrictPalindrome, err.Error())
			return config{}, err
		}
	}

	envMongoURI := os.Getenv("MONGO_URI")
	if *mongoURI == defaultMongoURI && envMongoURI != "" {
		*mongoURI = envMongoURI
	}

	return config{
		*httpAddr,
		*strictPalindrome,
		*mongoURI,
	}, nil
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
