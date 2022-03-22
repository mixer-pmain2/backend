package database

import (
	"database/sql"
	"fmt"

	_ "github.com/nakagami/firebirdsql"

	"pmain2/internal/apperror"
	"pmain2/internal/config"
)

type dbase interface {
	Close()
}

type DBase struct {
	DB *sql.DB
	dbase
}

func Connect() (*DBase, error) {
	conf, err := config.Create()
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("%s:%s@%s:%s/%s?encoding=WIN1251", conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbName)
	conn, err := sql.Open("firebirdsql", dsn)
	if err != nil {
		return nil, apperror.ErrDataBaseConnect
	}
	return &DBase{
		DB: conn,
	}, nil
}

func (db *DBase) Close() error {
	err := db.DB.Close()
	return err
}
