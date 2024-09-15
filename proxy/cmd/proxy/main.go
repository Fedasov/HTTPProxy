package main

import (
	"log"
	proxy "main/internal/pkg/http"
	"main/internal/server"
)

var (
	PROXY_PORT = "8080"
)

func main() {
	handler := proxy.CreateHandler()

	s := server.Init(PROXY_PORT, *handler)

	proto := "https"
	certPath := "certs/nck.crt"
	keyPath := "certs/cert.key"

	if proto == "http" {
		log.Fatal(s.ListenAndServe())
	} else {
		log.Fatal(s.ListenAndServeTLS(certPath, keyPath))
	}
}

/*
func main() {
		var pemPath string
		flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")
		var keyPath string
		flag.StringVar(&keyPath, "key", "server.key", "path to key file")
		var proto string
		flag.StringVar(&proto, "proto", "https", "Proxy protocol (http or https)")
		flag.Parse()
		if proto != "http" && proto != "https" {
			log.Fatal("Protocol must be either http or https")
		}

	if proto == "http" {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
	}
}
*/
