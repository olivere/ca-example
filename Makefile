.PHONY: server
server:
	go build -o srv ./server/...

.PHONY: client
client:
	go build -o cli ./client/...

start-server:
	./srv -addr=:10443 -ca=./certs/ca/ca.pem -key=./certs/servers/apiserver-key.pem -cert=./certs/servers/apiserver.pem

start-client:
	./cli -ca=./certs/ca/ca.pem -key=./certs/clients/admin@alt-f4.de-key.pem -cert=./certs/clients/admin@alt-f4.de.pem https://apiserver.go:10443
