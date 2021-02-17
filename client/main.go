package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	var (
		caFile         = flag.String("ca", "", "Certificate file, e.g. ca.pem")
		privateKeyFile = flag.String("key", "", "Private key file, e.g. apiserver-key.pem")
		publicKeyFile  = flag.String("cert", "", "Public key file, e.g. apiserver.pem")
	)
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	url := flag.Arg(0)
	if url == "" {
		log.Fatal("missing URL")
	}

	client, err := NewClient(
		WithCAFile(*caFile),
		WithClientCertFile(*privateKeyFile, *publicKeyFile),
		WithInsecureSkipVerify(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
