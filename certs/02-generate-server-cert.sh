#!/bin/sh
cfssl gencert \
    -ca=ca/ca.pem \
    -ca-key=ca/ca-key.pem \
    -config=ca/ca-config.json \
    -profile=server \
    servers/apiserver-csr.json | cfssljson -bare servers/apiserver
