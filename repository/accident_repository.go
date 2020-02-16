package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nelsonli33/go-parser/models"
)

type accidentRepository struct {
	db *sql.DB
}

func NewAccidentRepository(db *sql.DB) *accidentRepository {
	return &accidentRepository{db}
}

func (acc *accidentRepository) ClearAccidentTable() {
	stmt, err := acc.db.Query("TRUNCATE traffic_accident")

	if err != nil {
		panic(err)
	}

	defer stmt.Close()
}

func (acc *accidentRepository) InsertAccidents(accidents []models.TrafficAccident) {

	query := "INSERT INTO traffic_accident(date, death_count, injury_count, latitude, longitude) VALUES (?,?,?,?,?)"
	startTime := time.Now()
	tx, err := acc.db.Begin()
	if err != nil {
		panic(err)
	}
	for i, accident := range accidents {
		tx.Exec(query, accident.Date, accident.DeathCount, accident.InjuryCount, accident.Latitude, accident.Longitude)
		fmt.Printf("Data Index %d : %+v\n", i, accident)
	}
	tx.Commit()
	executeTime := time.Since(startTime)

	// 11675 32s without goroutine preparestatement
	// 11675 4s use the exec command in a loop inside a transaction.
	fmt.Printf("This sql execute: %s\n", executeTime)
	defer acc.db.Close()
}
