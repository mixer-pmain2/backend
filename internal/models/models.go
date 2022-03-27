package models

import (
	"strings"
	"time"

	"pmain2/pkg/logger"
	. "pmain2/pkg/utils"
)

var (
	INFO, _  = logger.New("dbase", logger.INFO)
	ERROR, _ = logger.New("dbase", logger.ERROR)
)

type SprDoct struct {
	Id       int    `json:"id"`
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

type Patient struct {
	Id         int    `json:"id"`
	Lname      string `json:"lname"`
	Fname      string `json:"fname"`
	Sname      string `json:"sname"`
	Bday       string `json:"bday"`
	Visibility int    `json:"visibility"`
	Sex        string `json:"sex"`
	Snils      string `json:"snils"`
	Address    string `json:"address"`
}

func (m *Patient) Serialize() error {
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
	m.Sex, err = ToUTF8(m.Sex)
	if err != err {
		return err
	}
	m.Address, err = ToUTF8(m.Address)
	if err != err {
		return err
	}
	m.Lname = strings.ReplaceAll(m.Lname, " ", "")
	m.Fname = strings.ReplaceAll(m.Fname, " ", "")
	m.Sname = strings.ReplaceAll(m.Sname, " ", "")
	m.Snils = strings.ReplaceAll(m.Snils, " ", "")
	m.Address = strings.ReplaceAll(m.Address, " ", "")
	bdayTime, err := time.Parse(time.RFC3339, m.Bday)
	if err != err {
		return err
	}
	m.Bday = ToDate(bdayTime)
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
