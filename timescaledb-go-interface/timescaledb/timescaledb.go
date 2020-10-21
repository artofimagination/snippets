package timescaledb

import (
	"database/sql"
	"log"
	"time"

	migrate "github.com/rubenv/sql-migrate"
	// Need to register postgres drivers with database/sql
	_ "github.com/lib/pq"
)

var DBAddress = "postgres://root:password@172.18.0.1:5432/data?sslmode=disable"

func BootstrapData() error {
	log.Println("Executing TimeScaleDB migration")

	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}
	log.Println("Getting migration files")

	db, err := sql.Open("postgres", DBAddress)
	if err != nil {
		return err
	}
	log.Println("DB connection open")

	retryCount := 5
	n := 0
	for retryCount > 0 {
		n, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
		if err != nil {
			retryCount--
			time.Sleep(1 * time.Second)
			log.Printf("Failed to execute migration %s. Retrying...\n", err.Error())
		} else {
			break
		}
	}

	if err != nil {
		return err
	}
	log.Printf("Applied %d migrations!\n", n)
	return nil
}

func ConnectData() (*sql.DB, error) {
	log.Println("Connecting to TimescaleDB")

	db, err := sql.Open("postgres", DBAddress)

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	return db, nil
}
