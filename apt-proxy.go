package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/lox/apt-proxy/proxy"
	"github.com/lox/httpcache"
)

const (
	defaultListen = "0.0.0.0:3142"
	defaultDir    = "./.aptcache"
)

var cachePatterns = proxy.CachePatternSlice{
	proxy.NewCachePattern(`deb$`, time.Hour*24*7),
	proxy.NewCachePattern(`udeb$`, time.Hour*24*7),
	proxy.NewCachePattern(`DiffIndex$`, time.Hour),
	proxy.NewCachePattern(`PackagesIndex$`, time.Hour),
	proxy.NewCachePattern(`Packages\.(bz2|gz|lzma)$`, time.Hour),
	proxy.NewCachePattern(`SourcesIndex$`, time.Hour),
	proxy.NewCachePattern(`Sources\.(bz2|gz|lzma)$`, time.Hour),
	proxy.NewCachePattern(`Release(\.gpg)?$`, time.Hour),
	proxy.NewCachePattern(`Translation-(en|fr)\.(gz|bz2|bzip2|lzma)$`, time.Hour),
	proxy.NewCachePattern(`Sources\.lzma$`, time.Hour),
}

var (
	version string
	listen  string
	dir     string
)

func init() {
	flag.StringVar(&listen, "listen", defaultListen, "the host and port to bind to")
	flag.StringVar(&dir, "cachedir", defaultDir, "the dir to store cache data in")
	flag.Parse()
}

func main() {
	log.Printf("running apt-proxy %s", version)

	cache, err := httpcache.NewDiskCache(dir)
	if err != nil {
		log.Fatal(err)
	}

	ap := proxy.NewAptProxy()
	ap.CachePatterns = cachePatterns

	log.Printf("proxy listening on https://%s", listen)
	log.Fatal(http.ListenAndServe(listen, &httpcache.Logger{
		Handler: httpcache.NewHandler(cache, ap.Handler()),
	}))
}
