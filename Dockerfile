FROM ubuntu:14.10

RUN apt-get update
RUN apt-get install -y golang git

ENV GOPATH /app
ADD . /app/src/github.com/lox/apt-proxy

WORKDIR /app/src/github.com/lox/apt-proxy
RUN go get

EXPOSE 8080
CMD ["go", "run", "/app/src/github.com/lox/apt-proxy/apt-proxy.go"]