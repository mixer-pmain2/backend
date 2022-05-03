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

type DiagM struct {
	Head         string `json:"head"`
	Diag         string `json:"diag"`
	Title        string `json:"title"`
	HaveChildren bool   `json:"haveChildren"`
}

func (m *SprModel) GetDiags(diag string) (*[]DiagM, error) {
	sql := fmt.Sprintf(`SELECT kod1, kod2, nam, uroven FROM diag1m where kod1 = '%s'`, diag)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	var data []DiagM
	for rows.Next() {
		row := DiagM{}
		err = rows.Scan(&row.Head, &row.Diag, &row.Title, &row.HaveChildren)
		if err != nil {
			return nil, err
		}
		row.Title, err = utils.ToUTF8(row.Title)
		row.Title = strings.Trim(row.Title, " ")
		row.Head = strings.Trim(row.Head, " ")
		row.Diag = strings.Trim(row.Diag, " ")
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	return &data, nil

}

type ServiceM struct {
	Param     string  `json:"param"`
	ParamI    int     `json:"paramI"`
	ParamD    float64 `json:"paramD"`
	ParamS    string  `json:"paramS"`
	Comment   string  `json:"comment"`
	DateStart string  `json:"dateStart"`
	DateEnd   string  `json:"dateEnd"`
}

func (m *SprModel) GetParams() (*[]ServiceM, error) {
	sql := fmt.Sprintf(`select param, PARAM_I, PARAM_D, PARAM_S, KOMMENT, DN, DK from servis`)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	var data []ServiceM
	for rows.Next() {
		row := ServiceM{}
		err = rows.Scan(&row.Param, &row.ParamI, &row.ParamD, &row.ParamS, &row.Comment, &row.DateStart, &row.DateEnd)
		if err != nil {
			return nil, err
		}
		row.ParamS, err = utils.ToUTF8(row.ParamS)
		if err != nil {
			return nil, err
		}
		row.Comment, err = utils.ToUTF8(row.Comment)
		if err != nil {
			return nil, err
		}
		row.ParamS = strings.Trim(row.ParamS, " ")
		row.Comment = strings.Trim(row.Comment, " ")
		row.Param = strings.Trim(row.Param, " ")
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	return &data, nil

}

func (m *SprModel) GetSprReason() (*map[string]string, error) {
	sql := fmt.Sprintf(`select kod1, na_me from spr_med where spr_nam = 'reg_reas1'
union 
select kod1, NA_ME from spr_med where spr_nam = 'exit_reas' and SUBSTR(kod1,1,1) = 'S' AND KOD2 = '1'`)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	data := make(map[string]string, 0)
	for rows.Next() {
		var row struct {
			name  string
			value string
		}
		err = rows.Scan(&row.name, &row.value)
		if err != nil {
			return nil, err
		}
		row.name = strings.Trim(row.name, " ")
		row.value = strings.Trim(row.value, " ")
		row.value, err = utils.ToUTF8(row.value)
		if err != nil {
			return nil, err
		}
		data[row.name] = row.value
	}
	return &data, nil

}
