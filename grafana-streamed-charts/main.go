package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type datapoint struct {
	PanelID int              `json:"panelid"` // Grafana panel id.
	RefID   string           `json:"refid"`   // Grafana panel query ref id.
	Values  map[string]int64 `json:"values"`  // Values associated with the row text. Timestamp is hardcoded every time.
}

type Config struct {
	Rows []string // contains the rows to display. Sent by grafana.
}

var configMap = make(map[int]map[string]Config)

func main() {
	// Grafana sends a request to initiate streaming in the form of http://localhost:8080?panelid=5&refid=A&data-rows=test1,test2
	http.HandleFunc("/", handler)
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

func handler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}

	// Parses the parameters.
	panelParam := r.URL.Query().Get("panelid")
	if panelParam == "" {
		panic("No panel id")
	}
	panelid, _ := strconv.Atoi(panelParam)

	rowsParam := r.URL.Query().Get("data-rows")
	if rowsParam == "" {
		panic("No rows")
	}

	refidParam := r.URL.Query().Get("refid")
	if refidParam == "" {
		panic("No ref id")
	}

	// Stores the row details for each query
	config := Config{
		Rows: strings.Split(rowsParam, ","),
	}
	queryConfig := make(map[string]Config)
	configMap[panelid] = queryConfig
	configMap[panelid][refidParam] = config
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		data := []datapoint{}
		offset := 0
		// For each query and for each row in the queries generate random values with 10 offset
		for panel, query := range configMap {
			for queryID, v := range query {
				dataRow := make(map[string]int64)
				dataRow["timestamp"] = time.Now().UnixNano() / 1000000
				for _, row := range v.Rows {
					random, _ := rand.Int(rand.Reader, big.NewInt(10))
					dataRow[row] = random.Int64() + int64(offset)
					offset += 10
				}
				currentPoint := &datapoint{
					RefID:   queryID,
					PanelID: panel,
					Values:  dataRow,
				}
				data = append(data, *currentPoint)
			}
		}

		j, _ := json.Marshal(data)
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		fmt.Fprintf(w, "%s\n", j)
		flusher.Flush() // Trigger "chunked" encoding and send a chunk...
	}
}
