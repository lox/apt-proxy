# Apt Proxy

A caching proxy specifically for apt package caching, also rewrites to the fastest local mirror.

Built because [apt-cacher-ng](https://www.unix-ag.uni-kl.de/~bloch/acng/) is unreliable.

## Running via Go

```
go install github.com/lox/apt-proxy
$GOBIN/apt-proxy
```