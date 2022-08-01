FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get -y update
RUN apt-get -y install vim build-essential golang-go openssh-client ca-certificates

RUN go install golang.org/x/tools/cmd/goimports@latest
ENV PATH="${PATH}:/root/go/bin"

RUN mkdir /usr/src/planning-poker
COPY go.mod go.sum *.go Makefile /usr/src/planning-poker

WORKDIR /usr/src/planning-poker
RUN make

ENTRYPOINT /usr/src/planning-poker/planning-poker
