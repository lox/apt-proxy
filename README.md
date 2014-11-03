# Apt Proxy

A caching proxy specifically for apt package caching, also rewrites to the fastest local mirror. Built as a tiny docker image for easy deployment.

Built because [apt-cacher-ng](https://www.unix-ag.uni-kl.de/~bloch/acng/) is unreliable.

## Running via Go

```bash
go install github.com/lox/apt-proxy
$GOBIN/apt-proxy
```

## Running in Docker for Development

```bash
docker build --rm --tag=apt-proxy-dev .
docker run -it --rm --publish=3142 --net host apt-proxy-dev
```

## Building in Docker for Release

```bash
docker build --rm --tag=apt-proxy-dev .
docker run -it --cidfile last-cid apt-proxy-dev ./build.sh
docker cp $(cat last-cid):/apt-proxy release/
docker build --tag=apt-proxy ./release
rm last-cid
```

## Running from Docker

```
docker run -it --rm --publish=3142 --net host lox24/apt-proxy
```
