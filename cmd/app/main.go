package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"github.com/google/uuid"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

const (
	transactionIDHeader  string = "X-Transaction-ID"
)

func transactionIdMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add(transactionIDHeader, uuid.New().String())
		next.ServeHTTP(w, r)
	})
}

func logTransactionTimeByIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w,r)
		elapsed := time.Since(start)
		log.Printf("Transaction %s took %s", r.Header.Get(transactionIDHeader), elapsed)
	})
}

func echo(w http.ResponseWriter, r *http.Request) {
	_, _ = io.Copy(w, r.Body)
}

func main() {
	var addr string

	flag.StringVar(&addr, "address", ":8080", "Address server will listen on")

	flag.Parse()

	http.Handle("/", transactionIdMiddleware(logTransactionTimeByIDMiddleware(http.HandlerFunc(echo))))

	log.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}