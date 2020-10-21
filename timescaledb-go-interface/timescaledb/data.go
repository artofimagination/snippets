package timescaledb

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type Data struct {
	CreatedAt time.Time       `json:"created_at"`
	ProjectID uuid.UUID       `json:"project_id"`
	RunSeqNo  int             `json:"run_seq_no"`
	Data      json.RawMessage `json:"data"`
}

// AddData will insert data into timescale db.
// Project ID is always generated in the user DB.
func AddData(projectID uuid.UUID, runSeqNo int, data interface{}) error {
	log.Println(projectID)
	query := "INSERT INTO project_data VALUES (NOW(), $1, $2, $3)"
	db, err := ConnectData()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, projectID, runSeqNo, data)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDataByProjectRun deletes all rows belonging to the selected run in the selected project.
func DeleteDataByProjectRun(projectID uuid.UUID, runSeqNo int) error {
	query := "DELETE FROM project_data WHERE project_id=$1 and run_seq_no=$2"
	db, err := ConnectData()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, projectID, runSeqNo)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDataByProject deletes all rows belonging to the projectID
func DeleteDataByProject(projectID uuid.UUID) error {
	query := "DELETE FROM project_data WHERE project_id=$1"
	db, err := ConnectData()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, projectID)
	if err != nil {
		return err
	}
	return nil
}

// GetDataByProjectRunChunk returns a chunk of data belonging to the specific run of the project with projectID.
// startTime defines the start time of the selection and itemCount refers to the number of rows to be returned after the startTime
func GetDataByProjectRunChunk(projectID uuid.UUID, runSeqNo int, startTime time.Time, itemCount int) (*[]Data, error) {
	log.Println(projectID)
	query := "SELECT * FROM project_data WHERE project_id = $1 AND run_seq_no = $2 and created_at > $3 limit $4"
	db, err := ConnectData()
	if err != nil {
		return &[]Data{}, err
	}
	defer db.Close()

	rows, err := db.Query(query, projectID, runSeqNo, startTime, itemCount)
	if err != nil {
		return &[]Data{}, err
	}

	dataList := []Data{}
	defer rows.Close()
	for rows.Next() {
		data := Data{}
		err = rows.Scan(&data.CreatedAt, &data.ProjectID, &data.RunSeqNo, &data.Data)
		if err != nil {
			return &[]Data{}, err
		}
		dataList = append(dataList, data)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return &[]Data{}, err
	}

	return &dataList, nil
}
