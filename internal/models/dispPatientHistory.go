package models

import (
	"fmt"
	"pmain2/pkg/utils"
	"strings"
	"time"
)

type HistoryVisit struct {
	Id       int    `json:"id"`
	Date     string `json:"date"`
	DockName string `json:"dockName"`
	Diag     string `json:"diag"`
	DiagS    string `json:"diagS"`
	Reason   string `json:"reason"`
	Where    string `json:"where"`
	Type     int    `json:"typeVisit"`
	Unit     int    `json:"unit"`
}

func (m *patientModel) HistoryVisits(id int) (*[]HistoryVisit, error) {
	var data []HistoryVisit
	sql := fmt.Sprintf(
		`select v_n, dat, dokf, diag, diag_t, prich, sost, m1, m2 from find_disp(%v)
     order by dat DESC`, id)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := HistoryVisit{}
		err = rows.Scan(&p.Id, &p.Date, &p.DockName, &p.Diag, &p.DiagS, &p.Reason, &p.Where, &p.Type, &p.Unit)
		dateVisit, _ := time.Parse(time.RFC3339, p.Date)
		p.Date = utils.ToDate(dateVisit)
		p.DockName, _ = utils.ToUTF8(p.DockName)
		p.Diag, _ = utils.ToUTF8(strings.Trim(p.Diag, " "))
		p.DiagS, _ = utils.ToUTF8(strings.Trim(p.DiagS, " "))
		p.Reason, _ = utils.ToUTF8(p.Reason)
		p.Where, _ = utils.ToUTF8(p.Where)
		if err != nil {
			return nil, err
		}
		data = append(data, p)
	}

	return &data, nil

}

type HistoryHospital struct {
	Id            int    `json:"id"`
	DateStart     string `json:"dateStart"`
	DateEnd       string `json:"dateEnd"`
	DiagStart     string `json:"diagStart"`
	DiagEnd       string `json:"diagEnd"`
	DiagStartS    string `json:"diagStartS"`
	DiagEndS      string `json:"diagEndS"`
	Otd           string `json:"otd"`
	HistoryNumber int    `json:"historyNumber"`
	Where         string `json:"where"`
}

func (m *patientModel) HistoryHospital(id int) (*[]HistoryHospital, error) {
	var data []HistoryHospital
	sql := fmt.Sprintf(`select datp, datv, dp, dv, diagp, diagv, otd, ni, we, nom_z from find_stac_otkaz(%v) order by datp DESC`, id)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := HistoryHospital{}
		err = rows.Scan(&p.DateStart, &p.DateEnd, &p.DiagStart, &p.DiagEnd, &p.DiagStartS, &p.DiagEndS, &p.Otd, &p.HistoryNumber, &p.Where, &p.Id)
		dateStart, _ := time.Parse(time.RFC3339, p.DateStart)
		p.DateStart = utils.ToDate(dateStart)
		dateEnd, _ := time.Parse(time.RFC3339, p.DateEnd)
		p.DateEnd = utils.ToDate(dateEnd)
		p.DiagStart = strings.Trim(p.DiagStart, " ")
		p.DiagEnd = strings.Trim(p.DiagEnd, " ")
		p.DiagStartS, _ = utils.ToUTF8(strings.Trim(p.DiagStartS, " "))
		p.DiagEndS, _ = utils.ToUTF8(strings.Trim(p.DiagEndS, " "))
		p.Where, _ = utils.ToUTF8(strings.Trim(p.Where, " "))
		if err != nil {
			return nil, err
		}
		data = append(data, p)
	}

	return &data, nil

}
