FROM ubuntu:14.10

RUN apt-get update
RUN apt-get install -y golang

RUN mkdir -p /go
ENV GOPATH /go
ENV GOBIN /go/bin
ADD . /go/src/github.com/lox/apt-proxy

EXPOSE 3142
WORKDIR /go/src/github.com/lox/apt-proxy
CMD ["go", "run", "/go/src/github.com/lox/apt-proxy/apt-proxy.go"]