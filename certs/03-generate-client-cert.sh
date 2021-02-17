#!/bin/sh
cfssl gencert \
    -ca=ca/ca.pem \
    -ca-key=ca/ca-key.pem  \
    -config=ca/ca-config.json \
    -profile=client \
    clients/admin@alt-f4.de-csr.json | cfssljson -bare clients/admin@alt-f4.de
