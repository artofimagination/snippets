package main

import (
	"flag"
	"net/http"
	"server-to-server-webrtc/webrtc"
	"time"
)

func sayHello() string {
	return "Hi! I am Server 1!"
}

func main() {
	webrtc.OfferAddr = flag.String("offer-address", "webrtc-client:8082", "Address that the Offer HTTP server is hosted on.")
	webrtc.AnswerAddr = flag.String("answer-address", "webrtc-server1:8080", "Address that the Answer HTTP server is hosted on.")
	webrtc.SendingFrequency = 500 * time.Millisecond
	candidateID := 1
	// Start Server
	webrtc.SetupServer(sayHello, candidateID)

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(*webrtc.AnswerAddr, nil))
	select {}
}
