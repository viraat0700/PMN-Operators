package main

import (
	// "fmt"
	// "io"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	values := strings.Split(r.URL.Path, "/")
	req := values[1]
	if req == "" {
		http.NotFound(w, r)
		return
	}
	log.Println("Server: ", req, " Client.")
	w.Write([]byte("Hello, Client !!.\n"))
	wg.Done()
}

func ServerService() {
	certPath := "../certs/server.crt"
	keyPath := "../certs/server.key"
	http.HandleFunc("/", HelloServer)
	err := http.ListenAndServeTLS(":2083", certPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("TlS Server has listing on port: 2083")
}

func ClientCallingToServer() {
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
	resp, err := client.Get("https://prometheus.nms.pmn-dev.wavelabs.in:2083/Hi")
	if err != nil {
		log.Println("Error is ", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error is ", err)
		return
	}
	defer resp.Body.Close()
	log.Println("Client cert are verified By Sever", resp.Status)
	log.Println("Respose sent By Server ", string(body))
	wg.Done()
}

func main() {
	wg.Add(2)
	go ServerService()
	go ClientCallingToServer()
	wg.Wait()
	//select {}
}
