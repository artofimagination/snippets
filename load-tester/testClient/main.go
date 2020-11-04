package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/template"
	"time"

	"github.com/artofimagination/grafana-json-streaming-datasource/streamer"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/gorilla/mux"
	"github.com/paulbellamy/ratecounter"
	"github.com/pkg/errors"
)

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

	empty := 0
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
	streamer.Origin = r.Header.Get("Origin")
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

	fillDataRow := func(row string, dataRow map[string]interface{}) {
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
	streamer.Configure(panelid, refidParam, rowsParam, startTime, endTime, datapoints, w)
	w.Header().Set("Access-Control-Allow-Origin", streamer.Origin)

	streamer.StreamData(panelid, refidParam, fillDataRow)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I am load tester!\n")
}
