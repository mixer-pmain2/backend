package models

import (
	"database/sql"
	"fmt"
	"pmain2/pkg/utils"
	"strings"
)

type SprModel struct {
	DB *sql.DB
}

func createSpr(db *sql.DB) *SprModel {
	return &SprModel{DB: db}
}

type PodrDict struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

func (m *SprModel) GetPodr() (*map[int]string, error) {
	sql := fmt.Sprintf(`SELECT MASKA1, NA_ME 
FROM SPR_PRAVA svn 
WHERE KOD1 = 1 AND MASKA1 > 0 and visible = 1`)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	var data = make(map[int]string, 32)
	for rows.Next() {
		row := PodrDict{}
		err = rows.Scan(&row.Code, &row.Name)
		if err != nil {
			return nil, err
		}
		row.Name, err = utils.ToUTF8(row.Name)
		row.Name = strings.Trim(row.Name, " ")
		if err != nil {
			return nil, err
		}
		//data = append(data, )
		data[row.Code] = row.Name
	}
	return &data, nil
}

type PravaDict struct {
	Unit int    `json:"unit"`
	Code int    `json:"code"`
	Name string `json:"name"`
}

func (m *SprModel) GetPrava() (*[]PravaDict, error) {
	sql := fmt.Sprintf(`SELECT MASKA1, MASKA2, NA_ME 
FROM SPR_PRAVA sp  
WHERE KOD1 = 2 and visible = 1
order by maska1, maska2`)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	var data = make([]PravaDict, 0)
	for rows.Next() {
		row := PravaDict{}
		err = rows.Scan(&row.Unit, &row.Code, &row.Name)
		if err != nil {
			return nil, err
		}
		row.Name, err = utils.ToUTF8(row.Name)
		row.Name = strings.Trim(row.Name, " ")
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	return &data, nil
}

type SprVisitD struct {
	Code int
	Name string
}

func (m *SprModel) GetSprVisit() (*map[int]string, error) {
	sql := fmt.Sprintf(`SELECT MASKA1, NA_ME 
FROM SPR_VISIT_N svn 
WHERE KOD1 = 3 AND MASKA1 > 0`)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	var data = make(map[int]string, 32)
	for rows.Next() {
		row := SprVisitD{}
		err = rows.Scan(&row.Code, &row.Name)
		if err != nil {
			return nil, err
		}
		row.Name, err = utils.ToUTF8(row.Name)
		row.Name = strings.Trim(row.Name, " ")
		if err != nil {
			return nil, err
		}
		data[row.Code] = row.Name
	}
	return &data, nil

}
