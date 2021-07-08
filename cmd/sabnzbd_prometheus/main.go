package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/triggity/sabnzbd_prometheus"
)

var (
	addrEnv      = "LISTEN_ADDRESS"
	sabUriEnv    = "SABNZBD_URI"
	sabApiKeyEnv = "SABNZBD_APIKEY"

	addr      = flag.String("listen-address", ":8081", fmt.Sprintf("The address to listen on for HTTP requests. Can also set via env %s", addrEnv))
	sabUri    = flag.String("sabnzbd-uri", "", fmt.Sprintf("the address for sabnzbd. Can also set via env %s", sabUriEnv))
	sabApiKey = flag.String("sabnzbd-apiKey", "", fmt.Sprintf("the apiKey for sabnzbd. Can also set via env %s", sabApiKeyEnv))
)

func getFlagValue(flagValue string, envKey string) string {
	envValue := os.Getenv(envKey)
	if envValue != "" {
		flagValue = envValue
	}
	return flagValue
}

func main() {

	flag.Parse()
	*addr = getFlagValue(*addr, addrEnv)
	*sabUri = getFlagValue(*sabUri, sabUriEnv)
	*sabApiKey = getFlagValue(*sabApiKey, sabApiKeyEnv)

	if *sabUri == "" {
		log.Fatalf("Must provide `-sabnzbd-uri` flag or %s environment variable\n", *sabUri)

	}
	if *sabApiKey == "" {
		log.Fatalf("Must provide `-sabnzbd-apiKey` flag or %s environment variable\n", *sabApiKey)
	}

	log.Printf("setting up sabnzbd client at %s\n", *sabUri)
	collector := sabnzbd_prometheus.NewSabNzbdCollector(*sabUri, *sabApiKey)
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("starting server, listening at %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
