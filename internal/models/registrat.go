package models

import (
	"database/sql"
)

type registratModel struct {
	DB *sql.DB
}

func createRegistrat(db *sql.DB) *registratModel {
	return &registratModel{DB: db}
}
