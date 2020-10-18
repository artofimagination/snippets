package webrtc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pion/logging"
	"github.com/pion/webrtc/v3"
)

// OfferAddr represents the address of the platform, that initiates
// webRTC with an offer.
var OfferAddr *string = nil

// AnswerAddr represents the product project address that is inactive
// until an offer arrives from the platform side.
var AnswerAddr *string = nil

// SendingFrequency determines how often a message is sent via webRTC.
var SendingFrequency time.Duration = 500 * time.Millisecond

var TurnServerAddress string = "turns:172.18.0.3:5349?transport=udp"
var StunServerAddress string = "stun:172.18.0.3:5349?transport=udp"
var TurnAuthCredential string = "default123secret"
var TurnAuthUserName string = "defaultUser"
var Certificates tls.Certificate = tls.Certificate{}

var config webrtc.Configuration = webrtc.Configuration{

	ICEServers: []webrtc.ICEServer{
		{
			URLs:           []string{TurnServerAddress, StunServerAddress},
			Username:       TurnAuthUserName,
			Credential:     TurnAuthCredential,
			CredentialType: webrtc.ICECredentialTypePassword,
		},
	},

	ICETransportPolicy: webrtc.ICETransportPolicyRelay,
}

// Everything below is the Pion WebRTC API! Thanks for using it ❤️.

// customLogger satisfies the interface logging.LeveledLogger
// a logger is created per subsystem in Pion, so you can have custom
// behavior per subsystem (ICE, DTLS, SCTP...)
type customLogger struct {
}

// Print all messages except trace
func (c customLogger) Trace(msg string) {
	fmt.Printf("customLogger Trace: %s\n", msg)
}
func (c customLogger) Tracef(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (c customLogger) Debug(msg string) { fmt.Printf("customLogger Debug: %s\n", msg) }
func (c customLogger) Debugf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}
func (c customLogger) Info(msg string) { fmt.Printf("customLogger Info: %s\n", msg) }
func (c customLogger) Infof(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}
func (c customLogger) Warn(msg string) { fmt.Printf("customLogger Warn: %s\n", msg) }
func (c customLogger) Warnf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}
func (c customLogger) Error(msg string) { fmt.Printf("customLogger Error: %s\n", msg) }
func (c customLogger) Errorf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// customLoggerFactory satisfies the interface logging.LoggerFactory
// This allows us to create different loggers per subsystem. So we can
// add custom behavior
type customLoggerFactory struct {
}

func (c customLoggerFactory) NewLogger(subsystem string) logging.LeveledLogger {
	fmt.Printf("Creating logger for %s \n", subsystem)
	return customLogger{}
}

func signalCandidate(client *http.Client, addr string, c *webrtc.ICECandidate) error {
	payload := []byte(c.ToJSON().Candidate)
	resp, err := client.Post(fmt.Sprintf("https://%s/candidate", addr), "application/json; charset=utf-8", bytes.NewReader(payload)) //nolint:noctx
	if err != nil {
		return err
	}

	if closeErr := resp.Body.Close(); closeErr != nil {
		return closeErr
	}

	return nil
}

// SetupProductSide sets up the webrtc connection on the product project side
// Each product project will have its own candidate id that is known for the platform side
// for identification purposes.
func SetupProductSide(client *http.Client, businessLogic func() string, candiateID int) {
	flag.Parse()

	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	s := webrtc.SettingEngine{
		LoggerFactory: customLoggerFactory{},
	}
	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))
	// Create a new RTCPeerConnection

	cert, err := webrtc.NewCertificate(Certificates.PrivateKey, *Certificates.Leaf)
	if err != nil {
		panic(err)
	}
	config.Certificates = []webrtc.Certificate{*cert}
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else if onICECandidateErr := signalCandidate(client, *OfferAddr, c); onICECandidateErr != nil {
			panic(onICECandidateErr)
		}
	})

	// A HTTP handler that allows the other Pion instance to send us ICE candidates
	// This allows us to add ICE candidates faster, we don't have to wait for STUN or TURN
	// candidates which may be slower
	candidatePath := fmt.Sprintf("/candidate")
	http.HandleFunc(candidatePath, func(w http.ResponseWriter, r *http.Request) {
		candidate, candidateErr := ioutil.ReadAll(r.Body)
		if candidateErr != nil {
			panic(candidateErr)
		}
		log.Println("Adding ICE Candidate")
		if candidateErr := peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: string(candidate)}); candidateErr != nil {
			panic(candidateErr)
		}
	})

	// A HTTP handler that processes a SessionDescription given to us from the other Pion process
	sdpPath := fmt.Sprintf("/sdp")
	http.HandleFunc(sdpPath, func(w http.ResponseWriter, r *http.Request) {
		sdp := webrtc.SessionDescription{}
		if err := json.NewDecoder(r.Body).Decode(&sdp); err != nil {
			panic(err)
		}

		if err := peerConnection.SetRemoteDescription(sdp); err != nil {
			panic(err)
		}

		// Create an answer to send to the other process
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			panic(err)
		}

		// Send our answer to the HTTP server listening in the other process
		payload, err := json.Marshal(answer)
		if err != nil {
			panic(err)
		}
		resp, err := client.Post(fmt.Sprintf("https://%s/sdp%d", *OfferAddr, candiateID), "application/json; charset=utf-8", bytes.NewReader(payload)) // nolint:noctx
		if err != nil {
			panic(err)
		} else if closeErr := resp.Body.Close(); closeErr != nil {
			panic(closeErr)
		}

		// Sets the LocalDescription, and starts our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			panic(err)
		}

		candidatesMux.Lock()
		for _, c := range pendingCandidates {
			onICECandidateErr := signalCandidate(client, *OfferAddr, c)
			if onICECandidateErr != nil {
				panic(onICECandidateErr)
			}
		}
		candidatesMux.Unlock()
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(SendingFrequency).C {
				message := businessLogic()
				fmt.Printf("Sending '%s'\n", message)

				// Send the message as text
				sendTextErr := d.SendText(message)
				if sendTextErr != nil {
					panic(sendTextErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})
}

// SetupPlatformSide initializes the webrtc connection and sends the initial offer
// to the product project candidates identified by candidate id.
func SetupPlatformSide(client *http.Client, businessLogic func() string, candiateID int) {
	flag.Parse()

	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	s := webrtc.SettingEngine{
		LoggerFactory: customLoggerFactory{},
	}
	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))

	cert, err := webrtc.NewCertificate(Certificates.PrivateKey, *Certificates.Leaf)
	if err != nil {
		panic(err)
	}
	config.Certificates = []webrtc.Certificate{*cert}
	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else if onICECandidateErr := signalCandidate(client, *AnswerAddr, c); err != nil {
			panic(onICECandidateErr)
		}
	})

	// A HTTP handler that allows the other Pion instance to send us ICE candidates
	// This allows us to add ICE candidates faster, we don't have to wait for STUN or TURN
	// candidates which may be slower
	candidatePath := fmt.Sprintf("/candidate%d", candiateID)
	http.HandleFunc(candidatePath, func(w http.ResponseWriter, r *http.Request) {
		candidate, candidateErr := ioutil.ReadAll(r.Body)
		if candidateErr != nil {
			panic(candidateErr)
		}
		if candidateErr := peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: string(candidate)}); candidateErr != nil {
			panic(candidateErr)
		}
	})

	// A HTTP handler that processes a SessionDescription given to us from the other Pion process
	sdpPath := fmt.Sprintf("/sdp%d", candiateID)
	http.HandleFunc(sdpPath, func(w http.ResponseWriter, r *http.Request) {
		sdp := webrtc.SessionDescription{}
		if sdpErr := json.NewDecoder(r.Body).Decode(&sdp); sdpErr != nil {
			panic(sdpErr)
		}

		if sdpErr := peerConnection.SetRemoteDescription(sdp); sdpErr != nil {
			panic(sdpErr)
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		for _, c := range pendingCandidates {
			if onICECandidateErr := signalCandidate(client, *AnswerAddr, c); onICECandidateErr != nil {
				panic(onICECandidateErr)
			}
		}
	})

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", dataChannel.Label(), dataChannel.ID())

		for range time.NewTicker(SendingFrequency).C {
			message := businessLogic()
			fmt.Printf("Sending '%s'\n", message)

			// Send the message as text
			sendTextErr := dataChannel.SendText(message)
			if sendTextErr != nil {
				panic(sendTextErr)
			}
		}
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	})

	// Create an offer to send to the other process
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	// Note: this will start the gathering of ICE candidates
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	// Send our offer to the HTTP server listening in the other process
	payload, err := json.Marshal(offer)
	if err != nil {
		panic(err)
	}
	resp, err := client.Post(fmt.Sprintf("https://%s/sdp", *AnswerAddr), "application/json; charset=utf-8", bytes.NewReader(payload)) // nolint:noctx
	if err != nil {
		panic(err)
	} else if err := resp.Body.Close(); err != nil {
		panic(err)
	}

}
