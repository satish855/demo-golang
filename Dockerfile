FROM golang:1.17.1
LABEL version ="0.1"
RUN export GOROOT=/usr/local/go/src
WORKDIR /usr/local/go/src/service
COPY . /usr/local/go/src/service
RUN pwd
RUN go build ./cmd/server
