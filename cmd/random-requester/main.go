package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"github.com/google/uuid"
)

var httpRequestMethods []string = []string{"DELETE", "GET", "PATCH", "POST", "PUT"}

func main() {
	var addr string

	flag.StringVar(&addr, "address", "", "Address to make requests of")

	flag.Parse()

	if addr == "" {
		log.Fatal(fmt.Printf("invalid address"))
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	client := &http.Client{}

	var buf *bytes.Buffer = bytes.NewBuffer([]byte{})
	var method string

	for {
		buf.Reset()
		buf.WriteString(uuid.New().String())

		method = httpRequestMethods[r.Intn(len(httpRequestMethods))]

		req, err := http.NewRequest(method, addr, buf)
		if err != nil {
			continue
		}

		start := time.Now()

		_, err = client.Do(req)
		if err != nil {
			log.Print(err)
		}

		elapsed := time.Since(start)
		log.Printf("%s request took %s", method, elapsed)

		time.Sleep(time.Duration(r.Intn(5)) * time.Second)
	}
}
