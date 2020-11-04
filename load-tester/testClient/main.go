package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"load-tester/mysqldb"

	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/gorilla/mux"
	"github.com/paulbellamy/ratecounter"
	"github.com/pkg/errors"
)

// Charts data representation.
type datapoint struct {
	PanelID int                    `json:"panelid"` // Grafana panel id.
	RefID   string                 `json:"refid"`   // Grafana panel query ref id.
	Values  map[string]interface{} `json:"values"`  // Values associated with the row text. Timestamp is hardcoded every time.
}

// QueryConfig contains ticker, response writer and series names
// for each series belonging to a query and for each query belonging to a panel.
type QueryConfig struct {
	Writer http.ResponseWriter
	Series []string
	Ticker *time.Ticker
}

var configMap = make(map[int]map[string]QueryConfig)
var origin = ""

// Placholders for measurements
var totalPrevCPU = uint64(0)
var usePrevCPU = uint64(0)
var avgCPU = ratecounter.NewAvgRateCounter(10 * time.Second)
var prevIOWait = uint64(0)
var avgIOWait = ratecounter.NewAvgRateCounter(10 * time.Second)

var insertExecCounter = ratecounter.NewAvgRateCounter(1 * time.Second)
var selectElapsed = time.Duration(0)
var insertRequestRateCount = ratecounter.NewRateCounter(1 * time.Second)
var failedInsertRequestCount = ratecounter.Counter(0)
var totalInsertElapsed = time.Time{}
var zeroTime = time.Time{}

func getCPUUsage() {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		panic("stat read fail")
	}

	totalCurr := stat.CPUStatAll.User + stat.CPUStatAll.System + stat.CPUStatAll.Idle
	useCurr := stat.CPUStatAll.User + stat.CPUStatAll.System
	ioWaitCurr := stat.CPUStatAll.IOWait

	diffIOWait := math.Abs(float64(ioWaitCurr - prevIOWait))
	diffTotal := math.Abs(float64(totalCurr - totalPrevCPU))
	diffUse := math.Abs(float64(useCurr - usePrevCPU))
	percentage := float32(0.0)
	if diffTotal != 0.0 {
		percentage = (float32(diffUse) * 100.0) / float32(diffTotal)
	}
	totalPrevCPU = totalCurr
	usePrevCPU = useCurr
	prevIOWait = ioWaitCurr

	avgIOWait.Incr(int64(diffIOWait))
	avgCPU.Incr(int64(percentage * 100))
}

func sendData(pID int, rID string) {
	dataRow := make(map[string]interface{})
	dataRow["timestamp"] = time.Now().UnixNano() / 1000000
	for _, row := range configMap[pID][rID].Series {
		switch row {
		case "cpu_load":
			getCPUUsage()
			dataRow[row] = float32(avgCPU.Rate() / 100.0)
		case "available_memory":
			dataRow[row] = getMemInfo()
		case "row_count":
			count, _ := mysqldb.GetUserCount()
			dataRow[row] = count
		case "insert_exec_time":
			dataRow[row] = insertExecCounter.Rate()
		case "select_exec_time":
			dataRow[row] = selectElapsed
		case "rate_count":
			dataRow[row] = insertRequestRateCount.Rate()
		case "failed_count":
			dataRow[row] = failedInsertRequestCount.Value()
		case "total_elapsed":
			dataRow[row] = time.Since(totalInsertElapsed).Nanoseconds()
		case "iowait":
			dataRow[row] = float32(1000000000 * avgIOWait.Rate() / 60.0) // convert Jiffies to nano second
		}
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

func getMemInfo() uint64 {
	info, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		panic("stat read fail")
	}
	return info.MemAvailable
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloServer)
	r.HandleFunc("/insert", insertIntoDB)
	r.HandleFunc("/select", selectFromDB)
	r.HandleFunc("/stream", streamHandler)
	r.HandleFunc("/show", showChart)

	if err := mysqldb.BootstrapSystem(); err != nil {
		log.Fatalf("System bootstrap failed. %s", errors.WithStack(err))
	}

	// Create Server and Route Handlers
	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	log.Println("Start shutting down")
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic("Failed to shut down server")
	}

	log.Println("Shutting down")
	os.Exit(0)
}

func showChart(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t := template.Must(template.ParseFiles(wd + "/charts/chart.html"))

	empty := datapoint{}
	err = t.ExecuteTemplate(w, "chart.html", empty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func insertIntoDB(w http.ResponseWriter, r *http.Request) {
	names, ok := r.URL.Query()["name"]
	if !ok || len(names[0]) < 1 {
		panic("Url Param 'name' is missing")
	}

	name := names[0]
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		panic("Url Param 'email' is missing")
	}

	email := emails[0]

	passwords, ok := r.URL.Query()["password"]
	if !ok || len(emails[0]) < 1 {
		panic("Url Param 'password' is missing")
	}

	password := passwords[0]

	start := time.Now()
	if totalInsertElapsed == zeroTime {
		totalInsertElapsed = start
	}
	err := mysqldb.AddUser(name, email, password)
	if err != nil {
		failedInsertRequestCount.Incr(1)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			panic("Failed to write error response")
		}
	}
	insertRequestRateCount.Incr(1)
	insertExecCounter.Incr(time.Since(start).Nanoseconds())
}

func selectFromDB(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	testValue := new(big.Int)
	fmt.Println(testValue.Binomial(1000, 10))

	selectElapsed = time.Since(start)
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

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I am load tester!\n")
}
