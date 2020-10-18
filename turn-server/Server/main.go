package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"test-server/webrtc"
	"time"
)

func sayHello() string {
	return "Hi! I am Server 1!"
}

func main() {
	webrtc.OfferAddr = flag.String("offer-address", "webrtc-client:8082", "Address that the Offer HTTP server is hosted on.")
	webrtc.AnswerAddr = flag.String("answer-address", "webrtc-server:8080", "Address that the Answer HTTP server is hosted on.")
	webrtc.SendingFrequency = 500 * time.Millisecond
	webrtc.TurnAuthUserName = "testUser"
	candidateID := 1

	cert, err := tls.LoadX509KeyPair("/etc/ssl/certs/cert.pem", "/etc/ssl/private/privkey.pem")
	if err != nil {
		log.Fatalf("failed to load server key pairs: %v", err)
	}

	pemFile, err := os.Open("/etc/ssl/certs/cert.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pemfileinfo, _ := pemFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(pemFile)
	_, err = buffer.Read(pembytes)
	block, _ := pem.Decode([]byte(pembytes))
	if block == nil {
		panic("failed to parse certificate PEM")

	}
	certPub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("failed to load system keychain: %v", err)
	}

	caCert, err := ioutil.ReadFile("/etc/ssl/certs/ca.pem")
	if err != nil {
		log.Fatalf("failed to read CA trust file rootCA.pem: %v", err)
	}

	ok := rootCAs.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatal("failed to load CA trust: bad PEM format?")
	}

	webrtc.Certificates = cert
	webrtc.Certificates.Leaf = certPub

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      rootCAs,
			},
		},
	}

	server := http.Server{
		Addr:    *webrtc.AnswerAddr,
		Handler: nil,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      rootCAs,
			//InsecureSkipVerify : true,
		},
	}

	go func() {
		panic(server.ListenAndServeTLS("", ""))
	}()
	// go func() {
	// 	panic(http.ListenAndServe(*webrtc.AnswerAddr, nil))
	// }()

	// Start Server
	webrtc.SetupProductSide(client, sayHello, candidateID)

	select {}
}
