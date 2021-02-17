#!/bin/sh
cfssl gencert -initca ca/ca-csr.json | cfssljson -bare ca/ca
