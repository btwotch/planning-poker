all: id_ecdsa planning-poker

id_ecdsa:
	ssh-keygen -q -N "" -t ecdsa -f id_ecdsa

planning-poker: *.go
	go fmt
	goimports -w .
	go mod tidy
	go build
	go test

