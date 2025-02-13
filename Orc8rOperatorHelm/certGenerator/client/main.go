package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)
	certPath := "../certs/client.crt"
	keyPath := "../certs/client.key"
	cert, err := ioutil.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Fatalf("Failed to read certificate file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certs, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Println(err)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certs},
			},
		},
	}

	resp, err := client.Get("https://localhost:2083/hello")
	if err != nil {
		fmt.Println("Error is ", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error is ", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("\n", resp.Status)
	fmt.Println(string(body))
}
