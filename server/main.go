package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

func main() {
	var (
		addr           = flag.String("addr", ":10443", "Listen address")
		caFile         = flag.String("ca", "", "Certificate file, e.g. ca.pem")
		privateKeyFile = flag.String("key", "", "Private key file, e.g. apiserver-key.pem")
		publicKeyFile  = flag.String("cert", "", "Public key file, e.g. apiserver.pem")
	)
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	s := newServer()

	errc := make(chan error, 1)

	go func() {
		h := http.Handler(s)
		errc <- serveTLS(h, *addr, *caFile, *privateKeyFile, *publicKeyFile)
	}()

	log.Fatal(<-errc)
}

func serveTLS(h http.Handler, addr string, caFile string, privateKeyFile, publicKeyFile string) error {
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(publicKeyFile, privateKeyFile)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		NextProtos:         []string{"h2", "http/1.1"},
		InsecureSkipVerify: true,
		// Client Auth
		ClientAuth: tls.VerifyClientCertIfGiven, // tls.RequireAndVerifyClientCert,
		ClientCAs:  pool,
	}

	tlsLn := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	server := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: h,
	}
	if err := http2.ConfigureServer(server, nil); err != nil {
		return err
	}
	log.Printf("Serving TLS at %s", tlsLn.Addr())
	return server.Serve(tlsLn)
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
