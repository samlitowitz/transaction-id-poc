package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"io"
	"log"
	"net/http"
	"github.com/google/uuid"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

const (
	transactionIDHeader string = "X-Transaction-ID"
)

type natsPubber struct {
	echoAddr string
	nc       *nats.Conn
	subj     string
}

func transactionIdMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add(transactionIDHeader, uuid.New().String())
		next.ServeHTTP(w, r)
	})
}

func (np *natsPubber) echoMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		client := &http.Client{}

		var buf *bytes.Buffer = bytes.NewBuffer([]byte{})

		req, err := http.NewRequest("GET", np.echoAddr, buf)
		if err != nil {
			return
		}
		req.Header.Add(transactionIDHeader, r.Header.Get(transactionIDHeader))

		resp, err := client.Do(req)
		if err != nil {
			log.Print(err)
			return
		}

		_, _ = io.Copy(w, resp.Body)
	})
}

func (np *natsPubber) logTransactionTimeByIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		msg := fmt.Sprintf("Transaction %s took %s", r.Header.Get(transactionIDHeader), elapsed)
		log.Print(msg)

		_ = np.nc.Publish(np.subj, []byte(msg))
		_ = np.nc.Flush()

		if err := np.nc.LastError(); err != nil {
			log.Fatal(err)
		}
	})
}

func echo(w http.ResponseWriter, r *http.Request) {
	_, _ = io.Copy(w, r.Body)
}

func main() {
	var addr, echoAddr, natsAddr, natsSubj string

	flag.StringVar(&addr, "address", ":8080", "Address server will listen on")
	flag.StringVar(&echoAddr, "echo-addr", "", "Echo service server")
	flag.StringVar(&natsAddr, "nats-addr", "", "Address of NATS server")
	flag.StringVar(&natsSubj, "nats-subj", "", "NATS channel")

	flag.Parse()

	if addr == "" {
		log.Fatal(fmt.Errorf("invalid address"))
	}
	if echoAddr == "" {
		log.Fatal(fmt.Errorf("invalid echo address"))
	}
	if natsAddr == "" {
		log.Fatal(fmt.Errorf("invalid NATS address"))
	}
	if natsSubj == "" {
		log.Fatal(fmt.Errorf("invalid NATS subject"))
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("Web Server")}

	// Connect to NATS
	nc, err := nats.Connect(natsAddr, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	np := &natsPubber{echoAddr: echoAddr, nc: nc, subj: natsSubj}

	http.Handle("/", transactionIdMiddleware(np.logTransactionTimeByIDMiddleware(np.echoMiddleWare(nil))))

	log.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
