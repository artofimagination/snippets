package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Charts data representation.
type datapoint struct {
	PanelID int                    `json:"panelid"` // Grafana panel id.
	RefID   string                 `json:"refid"`   // Grafana panel query ref id.
	Values  map[string]interface{} `json:"values"`  // Values associated with the row text. Timestamp is hardcoded every time.
}

type QueryConfig struct {
	Writer http.ResponseWriter
	Series []string
	Ticker *time.Ticker
}

var configMap = make(map[int]map[string]QueryConfig)
var origin = ""

func main() {
	// Grafana sends a request to initiate streaming in the form of http://localhost:8080?panelid=5&refid=A&data-rows=test1,test2
	http.HandleFunc("/", streamHandler)
	http.HandleFunc("/show", showChart)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func showChart(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t := template.Must(template.ParseFiles(wd + "/chart.html"))

	empty := datapoint{}
	err = t.ExecuteTemplate(w, "chart.html", empty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sendData(pID int, rID string) {
	dataRow := make(map[string]interface{})
	dataRow["timestamp"] = time.Now().UnixNano() / 1000000
	offset := 0
	for _, row := range configMap[pID][rID].Series {
		switch row {
		case "cpu_load":
			random, _ := rand.Int(rand.Reader, big.NewInt(10))
			dataRow[row] = random.Int64() + int64(offset)
		case "available_memory":
			random, _ := rand.Int(rand.Reader, big.NewInt(10))
			dataRow[row] = random.Int64() + int64(offset)
		case "row_count":
			random, _ := rand.Int(rand.Reader, big.NewInt(10))
			dataRow[row] = random.Int64() + int64(offset)
		}
		offset += 10
	}
	currentPoint := &datapoint{
		RefID:   rID,
		PanelID: pID,
		Values:  dataRow,
	}
	j, _ := json.Marshal(currentPoint)
	configMap[pID][rID].Writer.Header().Set("Access-Control-Allow-Origin", origin)
	fmt.Fprintf(configMap[pID][rID].Writer, "%s\n", j)
	configMap[pID][rID].Writer.(http.Flusher).Flush() // Trigger "chunked" encoding and s
}

func streamData(pID int, rID string) {
	log.Println("Start streaming")
	log.Println(configMap)
	defer configMap[pID][rID].Ticker.Stop()
	for ; true; <-configMap[pID][rID].Ticker.C {
		sendData(pID, rID)
	}
}

// Whenever a grafana query request arrives this handler is called.
// Based on the parsed parameters the configMap is extended
// Structure of config map:
//			Panel -> queries/panel -> series/queries
// At the moment each query has its own ticker channel setup based on the sampling time (calculated based on the grafana data points config)
// Each series is represented by a string that can be configured in grafana through the dataText field.
// Once the ticker is set, the respective query results will be streamed back (periodical http flush) to the grafana server.
func streamHandler(w http.ResponseWriter, r *http.Request) {
	origin = r.Header.Get("Origin")
	// Parses the parameters.
	panelParam := r.URL.Query().Get("panelid")
	if panelParam == "" {
		panic("No panel id")
	}
	panelid, _ := strconv.Atoi(panelParam)

	startTimeParam := r.URL.Query().Get("start")
	if startTimeParam == "" {
		panic("No start time")
	}
	startTime, _ := strconv.Atoi(startTimeParam)

	endTimeParam := r.URL.Query().Get("end")
	if endTimeParam == "" {
		panic("No end time")
	}
	endTime, _ := strconv.Atoi(endTimeParam)

	datapointsParam := r.URL.Query().Get("datapoints")
	if datapointsParam == "" {
		panic("No datapoints")
	}
	datapoints, _ := strconv.Atoi(datapointsParam)

	rowsParam := r.URL.Query().Get("data-rows")
	if rowsParam == "" {
		panic("No rows")
	}

	refidParam := r.URL.Query().Get("refid")
	if refidParam == "" {
		panic("No ref id")
	}

	// Sets the sampling time based on the allowed grafana data points
	samplingTimeMs := (endTime - startTime) / datapoints
	if _, ok := configMap[panelid]; !ok {
		configMap[panelid] = make(map[string]QueryConfig)
	}

	// Creates the config map
	configMap[panelid][refidParam] = QueryConfig{
		Series: strings.Split(rowsParam, ","),
		Ticker: time.NewTicker(time.Duration(samplingTimeMs) * time.Millisecond),
		Writer: w,
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)

	streamData(panelid, refidParam)
}
