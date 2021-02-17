# CA infrastructure

This is an example of using TLS Client Authentication with certificates.

## Prerequisites

We use the [cfssl Suite](https://github.com/cloudflare/cfssl) to generate
and configure certificates.

Also, make sure that `apiserver.go` resolves to `127.0.0.1` by e.g. adding
it to `/etc/hosts`.

## Create certificate infrastructure

### Create the CA

Step 1 is to set up your CA infrastructure. First, create the CA:

```
$ cd certs
$ ./01-generate-root-ca.sh
2018/10/11 11:40:17 [INFO] generating a new CA key and certificate from CSR
2018/10/11 11:40:17 [INFO] generate received request
2018/10/11 11:40:17 [INFO] received CSR
2018/10/11 11:40:17 [INFO] generating key: rsa-2048
2018/10/11 11:40:17 [INFO] encoded CSR
2018/10/11 11:40:17 [INFO] signed certificate with serial number 123123123123...
$ ls ca/
ca-config.json	ca-csr.json	ca-key.pem	ca.csr		ca.pem
```

The `ca/ca-key.pem` file is the private key while `ca/ca.pem` is the public
key file for your CA.

### Create the server certificate

Step 2 is to create a public key pair for the server-side:

```
$ ./02-generate-server-cert.sh
2018/10/11 11:42:05 [INFO] generate received request
2018/10/11 11:42:05 [INFO] received CSR
2018/10/11 11:42:05 [INFO] generating key: rsa-2048
2018/10/11 11:42:05 [INFO] encoded CSR
2018/10/11 11:42:05 [INFO] signed certificate with serial number 123123...
$ ls servers
apiserver-csr.json	apiserver-key.pem	apiserver.csr		apiserver.pem
```

Now we have a certificate for `apiserver.go`. Again,
`servers/apiserver-key.pem` holds the private key while
`servers/apiserver.pem` holds the public key.

### Create the client certificates

Step 3 creates client certificates. You might need as many as there are
users calling the server. Here's how to set up a client with the common
name of `admin@alt-f4.de`:

```
$ ./03-generate-client-cert.sh
2018/10/11 11:44:52 [INFO] generate received request
2018/10/11 11:44:52 [INFO] received CSR
2018/10/11 11:44:52 [INFO] generating key: rsa-2048
2018/10/11 11:44:52 [INFO] encoded CSR
2018/10/11 11:44:52 [INFO] signed certificate with serial number 123123123...
$ ls clients
admin@alt-f4.de-csr.json	admin@alt-f4.de-key.pem	admin@alt-f4.de.csr		admin@alt-f4.de.pem
```

## Starting the server

The server is configured to use `tls.VerifyClientCertIfGiven` by default,
meaning: You can give me a certificate and if you do, I'll verify it. If
you don't, I'll proceed the request anyway. You can change that in
[`server/main.go`](https://github.com/olivere/ca-example/blob/master/server/main.go).

Build and run the server:

```
$ make server start-server
go build -o srv ./server/...
./srv -addr=:10443 -ca=./certs/ca/ca.pem -key=./certs/servers/apiserver-key.pem -cert=./certs/servers/apiserver.pem
2018/10/11 11:50:08 Serving TLS at [::]:10443
```

## Starting the client

To run the client, do:

```
$ make client start-client
go build -o cli ./client/...
./cli -ca=./certs/ca/ca.pem -key=./certs/clients/admin@alt-f4.de-key.pem -cert=./certs/clients/admin@alt-f4.de.pem https://apiserver.go:10443
Authenticated as admin@alt-f4.de
```

To see what happens if you don't specify a valid certificate, do:

```
$ curl -k -I https://apiserver.go:10443
HTTP/2 401
content-type: text/plain; charset=utf-8
x-content-type-options: nosniff
content-length: 35
date: Thu, 11 Oct 2018 09:52:42 GMT

```

# License

Copyright (c) 2021 Oliver Eilhard. All rights reserved.
