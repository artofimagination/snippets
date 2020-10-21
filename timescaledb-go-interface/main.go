package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"timescaledb-go-interface/jsonutils"
	"timescaledb-go-interface/timescaledb"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am Server!")
}

func insertData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Inserting data")
	project, err := uuid.NewUUID()
	if err != nil {
		fmt.Fprintln(w, "Cannot generate 'project' UUID")
		return
	}

	projects, ok := r.URL.Query()["project"]
	if ok && len(projects) == 1 {
		project = uuid.MustParse(projects[0])
	}

	seqNos, ok := r.URL.Query()["seqNo"]
	if !ok || len(seqNos[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'seqNo' is missing")
		return
	}

	seqNo := seqNos[0]

	dataList, ok := r.URL.Query()["data"]
	if !ok || len(dataList[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'data' is missing")
		return
	}

	data := dataList[0]

	seqNoInt, err := strconv.Atoi(seqNo)
	if err != nil {
		fmt.Fprintln(w, "Failed to convert 'seqNo' string to int")
		return
	}

	err = timescaledb.AddData(project, seqNoInt, data)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
}

func deleteDataByProjectRun(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Deleting data by run")
	projects, ok := r.URL.Query()["project"]
	if !ok || len(projects[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'project' is missing")
		return
	}
	project := projects[0]

	seqNos, ok := r.URL.Query()["seqNo"]
	if !ok || len(seqNos[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'seqNo' is missing")
		return
	}

	seqNo := seqNos[0]

	seqNoInt, err := strconv.Atoi(seqNo)
	if err != nil {
		fmt.Fprintln(w, "Failed to convert 'seqNo' string to int")
		return
	}

	err = timescaledb.DeleteDataByProjectRun(uuid.MustParse(project), seqNoInt)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
}

func deleteDataByProject(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Deleting data by project")
	projects, ok := r.URL.Query()["project"]
	if !ok || len(projects[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'project' is missing")
		return
	}
	project := projects[0]

	err := timescaledb.DeleteDataByProject(uuid.MustParse(project))
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
}

func getDataByProjectRun(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Getting data")
	projects, ok := r.URL.Query()["project"]
	if !ok || len(projects[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'project' is missing")
		return
	}

	project := projects[0]

	seqNos, ok := r.URL.Query()["seqNo"]
	if !ok || len(seqNos[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'seqNo' is missing")
		return
	}

	seqNo := seqNos[0]
	seqNoInt, err := strconv.Atoi(seqNo)
	if err != nil {
		fmt.Fprintln(w, "Failed to convert 'seqNo' string to int")
		return
	}

	chunkSizes, ok := r.URL.Query()["chunk"]
	if !ok || len(chunkSizes[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'chunk' is missing")
		return
	}

	chunk := chunkSizes[0]
	chunkInt, err := strconv.Atoi(chunk)
	if err != nil {
		fmt.Fprintln(w, "Failed to convert 'chunk' string to int")
		return
	}

	dates, ok := r.URL.Query()["date"]
	if !ok || len(dates[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'date' is missing")
		return
	}

	date := dates[0]
	startTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	data, err := timescaledb.GetDataByProjectRunChunk(uuid.MustParse(project), seqNoInt, startTime, chunkInt)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, data)
		for _, element := range *data {
			fmt.Fprintln(w, element)
			dataJSON, err := jsonutils.ProcessJSON(element.Data)
			if err != nil {
				fmt.Fprintln(w, err.Error())
				return
			}
			fmt.Fprintln(w, dataJSON)
		}
	}
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/insert", insertData)
	http.HandleFunc("/get-by-project-run", getDataByProjectRun)
	http.HandleFunc("/delete-by-project-run", deleteDataByProjectRun)
	http.HandleFunc("/delete-by-project", deleteDataByProject)

	if err := timescaledb.BootstrapData(); err != nil {
		log.Fatalf("Data bootstrap failed. %s", errors.WithStack(err))
	}

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(":8080", nil))
}
