package models

import (
	"database/sql"
	"pmain2/pkg/logger"
	. "pmain2/pkg/utils"
	"strings"
)

var (
	INFO, _  = logger.New("dbase", logger.INFO)
	ERROR, _ = logger.New("dbase", logger.ERROR)

	db *sql.DB

	Model *models
)

type models struct {
	Patient        patientModel
	Doctor         doctorModel
	Registrat      registratModel
	Spr            SprModel
	User           userModel
	Visit          VisitModel
	Administration administrationModel
}

func Init(dbase *sql.DB) *models {
	db = dbase
	INFO.Println("Init models")
	Model = &models{
		Patient:   *createPatient(),
		Doctor:    *createDoctor(),
		Registrat: *createRegistrat(),
		Spr:       *createSpr(),
		User:      *createUser(),
		Visit:     *createVisit(),
	}
	return Model
}

func (m *models) CreateTx() (error, *sql.Tx) {

	tx, err := db.Begin()
	if err != nil {
		return err, nil
	}
	return nil, tx
}

type SprDoct struct {
	Id       int64  `json:"id"`
	Lname    string `json:"lname"`
	Fname    string `json:"fname"`
	Sname    string `json:"sname"`
	Password string `json:"-"`
}

func (m *SprDoct) Serialize() error {
	var err error
	m.Lname, err = ToUTF8(m.Lname)
	if err != err {
		return err
	}
	m.Fname, err = ToUTF8(m.Fname)
	if err != err {
		return err
	}
	m.Sname, err = ToUTF8(m.Sname)
	if err != err {
		return err
	}
	m.Lname = strings.ReplaceAll(m.Lname, " ", "")
	m.Fname = strings.ReplaceAll(m.Fname, " ", "")
	m.Sname = strings.ReplaceAll(m.Sname, " ", "")
	return nil
}

type Visit struct {
	Id        int    `json:"id"`
	PatientId int    `json:"patientId"`
	Date      string `json:"date"`
	DockId    int    `json:"doctId"`
	Diagnose  string `json:"diagnose"`
	Type      int    `json:"type"`
	Pord      int    `json:"pord"`
	Home      bool   `json:"home"`
}

type Registrat struct {
	Id        int
	PatientId int
	Uch       int
	RegDate   string
	DockId    int
	Reason    string
	Category  int
	Diagnose  string
}
