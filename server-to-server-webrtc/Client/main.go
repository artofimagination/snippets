package main

import (
	"flag"
	"log"
	"net/http"
	"server-to-server-webrtc/webrtc"
	"time"
)

func sayHello() string {
	return "Hi! I am Client!"
}

func main() {
	webrtc.OfferAddr = flag.String("offer-address", "webrtc-client:8082", "Address that the Offer HTTP server is hosted on.")
	webrtc.AnswerAddr = flag.String("answer-address1", "webrtc-server1:8080", "Address that the Answer HTTP server is hosted on.")
	webrtc.SendingFrequency = 500 * time.Millisecond
	candidateID := 1
	// Start HTTP server that accepts requests from the answer process
	go func() { panic(http.ListenAndServe(*webrtc.OfferAddr, nil)) }()
	// Start Server
	webrtc.SetupPlatformSide(sayHello, candidateID)
	log.Println("First test")

	webrtc.AnswerAddr = flag.String("answer-address2", "webrtc-server2:8081", "Address that the Answer HTTP server is hosted on.")
	webrtc.SendingFrequency = 500 * time.Millisecond
	candidateID = 2
	webrtc.SetupPlatformSide(sayHello, candidateID)
	log.Println("Second test")
	select {}
}
