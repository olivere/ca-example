package main

import (
	"fmt"
	"net/http"
)

var (
	_ http.Handler = (*server)(nil)
)

type server struct {
}

func newServer() *server {
	return &server{}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		// HTTP -> HTTPS
		if r.Host == "" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusFound)
		return
	}

	if len(r.TLS.PeerCertificates) == 0 {
		http.Error(w, "not authenticated by a certificate", http.StatusUnauthorized)
		return
	}

	cert := r.TLS.PeerCertificates[0]
	fmt.Fprintf(w, "Authenticated as %s\n", cert.Subject.CommonName)
}
