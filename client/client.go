package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

func NewClient(opts ...ClientOption) (*http.Client, error) {
	var cli client
	for _, o := range opts {
		if err := o(&cli); err != nil {
			return nil, err
		}
	}

	c := &http.Client{}
	if cli.config != nil {
		cli.config.BuildNameToCertificate()
		tr := &http.Transport{
			TLSClientConfig: cli.config,
		}
		c.Transport = tr
	}
	return c, nil
}

type client struct {
	config *tls.Config
}

type ClientOption func(*client) error

func WithCAFile(caFile string) ClientOption {
	return func(c *client) error {
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caCert)

		if c.config == nil {
			c.config = &tls.Config{}
		}
		c.config.RootCAs = pool

		return nil
	}
}

func WithClientCertFile(privateKeyFile, publicKeyFile string) ClientOption {
	return func(c *client) error {
		cert, err := tls.LoadX509KeyPair(publicKeyFile, privateKeyFile)
		if err != nil {
			return err
		}
		certs := []tls.Certificate{cert}
		if c.config == nil {
			c.config = &tls.Config{}
		}
		c.config.Certificates = certs
		return nil
	}
}

func WithInsecureSkipVerify(skip bool) ClientOption {
	return func(c *client) error {
		if c.config == nil {
			c.config = &tls.Config{}
		}
		c.config.InsecureSkipVerify = skip
		return nil
	}
}
