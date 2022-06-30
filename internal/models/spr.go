package models

import (
	"database/sql"
	"fmt"
	"pmain2/internal/types"
	"pmain2/pkg/utils"
	"strconv"
	"strings"
	"time"
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
	INFO.Println(sql)
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
	INFO.Println(sql)
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
	INFO.Println(sql)
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
	INFO.Println(sql)
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
	INFO.Println(sql)
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
	INFO.Println(sql)
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

func (m *SprModel) IsClosedSection(section int) (bool, error) {
	sqlQuery := fmt.Sprintf(`select First 1 CLOSED_DAT from closed_uch where uch = %v`, section)
	INFO.Println(sqlQuery)
	row := m.DB.QueryRow(sqlQuery)

	isClose := false
	var _closeDate string
	err := row.Scan(&_closeDate)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if _closeDate == "" {
		return false, nil
	}
	closeDate, err := time.Parse(time.RFC3339, _closeDate)
	if err != nil {
		return false, err
	}
	curDate := time.Now()
	if curDate.Sub(closeDate) > 0 {
		isClose = true
	}

	return isClose, nil

}

func (m *SprModel) GetSprInvalidKind() (*map[string]string, error) {
	sql := fmt.Sprintf(`select kod2, na_me from spr_visit_n
where kod1 = 6 and kod2 > 0 and visible = 1
order by na_me`)
	INFO.Println(sql)
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

func (m *SprModel) GetSprInvalidChildAnomaly() (*map[string]string, error) {
	sql := fmt.Sprintf(`select kod2, na_me from spr_visit_n
where kod1 = 14 and kod2 > 0 and visible = 1
order by na_me`)
	INFO.Println(sql)
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

func (m *SprModel) GetSprInvalidChildLimit() (*map[string]string, error) {
	sql := fmt.Sprintf(`select kod2, na_me from spr_visit_n
where kod1 = 12 and kod2 > 0`)
	INFO.Println(sql)
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

func (m *SprModel) GetSprInvalidReason() (*map[string]string, error) {
	sql := fmt.Sprintf(`select kod2, na_me from spr_visit_n
where kod1 = 4 and kod2 > 0 and visible = 1`)
	INFO.Println(sql)
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

func (m *SprModel) GetSprCustodyWho() (*map[string]string, error) {
	sql := fmt.Sprintf(`select kod1, NA_ME from spr_med
where spr_nam = 'care'`)
	INFO.Println(sql)
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

func (m *SprModel) FindRepublic(find *types.Find) (*[]types.Spr, error) {
	sql := fmt.Sprintf(`select kod1, na_me from spr_adress
where spr_nam = 'republic' and na_me >= ?
order by na_me`)
	INFO.Println(sql)
	stmt, err := m.DB.Prepare(sql)
	rows, err := stmt.Query(find.Name)
	if err != nil {
		return nil, err
	}
	data := make([]types.Spr, 0)
	for rows.Next() {
		var row types.Spr
		err = rows.Scan(&row.Code, &row.Name)
		if err != nil {
			return nil, err
		}
		row.Code, _ = utils.ToUTF8(strings.Trim(row.Code, " "))
		row.Name, _ = utils.ToUTF8(strings.Trim(row.Name, " "))
		data = append(data, row)
	}
	return &data, nil

}

func (m *SprModel) FindRegion(find *types.Find) (*[]types.Spr, error) {
	sql := fmt.Sprintf(`select kod1, na_me from spr_adress
where spr_nam = 'region' and na_me >= ?
order by na_me
`)
	INFO.Println(sql)
	stmt, err := m.DB.Prepare(sql)
	rows, err := stmt.Query(find.Name)
	if err != nil {
		return nil, err
	}
	data := make([]types.Spr, 0)
	for rows.Next() {
		var row types.Spr
		err = rows.Scan(&row.Code, &row.Name)
		if err != nil {
			return nil, err
		}
		row.Code, _ = utils.ToUTF8(strings.Trim(row.Code, " "))
		row.Name, _ = utils.ToUTF8(strings.Trim(row.Name, " "))
		data = append(data, row)
	}
	return &data, nil

}

func (m *SprModel) FindDistrict(find *types.Find) (*[]types.Spr, error) {
	sql := fmt.Sprintf(`select kod1, na_me from spr_adress
where spr_nam = 'district' and na_me >= ?
order by na_me`)
	INFO.Println(sql)
	stmt, err := m.DB.Prepare(sql)
	rows, err := stmt.Query(find.Name)
	if err != nil {
		return nil, err
	}
	data := make([]types.Spr, 0)
	for rows.Next() {
		var row types.Spr
		err = rows.Scan(&row.Code, &row.Name)
		if err != nil {
			return nil, err
		}
		row.Code, _ = utils.ToUTF8(strings.Trim(row.Code, " "))
		row.Name, _ = utils.ToUTF8(strings.Trim(row.Name, " "))
		data = append(data, row)
	}
	return &data, nil

}

func (m *SprModel) FindArea(find *types.Find) (*[]types.Spr, error) {
	sql := fmt.Sprintf(`select kod1, na_me, kod2 from spr_adress
where spr_nam = 'pop_area' and na_me >= ?
order by na_me`)
	INFO.Println(sql)
	stmt, err := m.DB.Prepare(sql)
	rows, err := stmt.Query(find.Name)
	if err != nil {
		return nil, err
	}
	data := make([]types.Spr, 0)
	for rows.Next() {
		var row types.Spr
		err = rows.Scan(&row.Code, &row.Name, &row.Param)
		if err != nil {
			return nil, err
		}
		row.Code, _ = utils.ToUTF8(strings.Trim(row.Code, " "))
		row.Name, _ = utils.ToUTF8(strings.Trim(row.Name, " "))
		row.Param, _ = utils.ToUTF8(strings.Trim(row.Param, " "))
		data = append(data, row)
	}
	return &data, nil

}

func (m *SprModel) FindStreet(find *types.Find) (*[]types.Spr, error) {
	sql := fmt.Sprintf(`select kod1, na_me from spr_adress
where spr_nam = 'street' and na_me >= ?
order by na_me`)
	INFO.Println(sql)
	stmt, err := m.DB.Prepare(sql)
	rows, err := stmt.Query(find.Name)
	if err != nil {
		return nil, err
	}
	data := make([]types.Spr, 0)
	for rows.Next() {
		var row types.Spr
		err = rows.Scan(&row.Code, &row.Name)
		if err != nil {
			return nil, err
		}
		row.Code, _ = utils.ToUTF8(strings.Trim(row.Code, " "))
		row.Name, _ = utils.ToUTF8(strings.Trim(row.Name, " "))
		data = append(data, row)
	}
	return &data, nil

}

func (m *SprModel) FindSections(find *types.FindI) (*[]types.SprUchN, error) {
	sql := fmt.Sprintf(`SELECT nom_z, uch, KOMMENT, PLAN_P, CHAS, SPEC, disp FROM SPR_UCH_N sun`)
	if find.Name > 0 {
		sql = fmt.Sprintf(`%s where bin_and(disp,?) = ? and uch >= 10`, sql)
	}
	INFO.Println(sql)
	stmt, err := m.DB.Prepare(sql)
	rows, err := stmt.Query(find.Name, find.Name)
	if err != nil {
		return nil, err
	}
	data := make([]types.SprUchN, 0)
	for rows.Next() {
		var row types.SprUchN
		err = rows.Scan(&row.Id, &row.Section, &row.Name, &row.Plan, &row.Hour, &row.Spec, &row.Unit)
		if err != nil {
			return nil, err
		}
		row.Name, _ = utils.ToUTF8(strings.Trim(row.Name, " "))
		row.Spec, _ = utils.ToUTF8(strings.Trim(row.Spec, " "))
		data = append(data, row)
	}
	return &data, nil

}

func (m *SprModel) FindSectionDoctor(find *types.FindI) (*[]types.LocationDoctor, error) {
	sql := fmt.Sprintf(`SELECT du.uch, du.DOCK, sd.FIO, sd.IM, sd.OT, du.PODRAZ 
FROM DOCK_UCH du
LEFT JOIN spr_doct sd ON sd.KOD_DOCK = du.DOCK 
WHERE du.PRIZ = 1
AND PODRAZ = ?
ORDER BY uch`)
	INFO.Println(sql)
	stmt, err := m.DB.Prepare(sql)
	rows, err := stmt.Query(find.Name)
	if err != nil {
		return nil, err
	}
	data := make([]types.LocationDoctor, 0)
	for rows.Next() {
		var row types.LocationDoctor
		id := ""
		err = rows.Scan(&row.Section, &id, &row.Lname, &row.Fname, &row.Sname, &row.Unit)
		if err != nil {
			return nil, err
		}
		row.DoctId, _ = strconv.Atoi(strings.Trim(id, " "))
		row.Lname, _ = utils.ToUTF8(strings.Trim(row.Lname, " "))
		row.Fname, _ = utils.ToUTF8(strings.Trim(row.Fname, " "))
		row.Sname, _ = utils.ToUTF8(strings.Trim(row.Sname, " "))
		data = append(data, row)
	}
	return &data, nil
}

func (m *SprModel) GetDoctors(find *types.FindI) (*[]types.Doctor, error) {
	sql := fmt.Sprintf(`SELECT KOD_DOCK_I, FIO, IM, OT, prava, z152
FROM SPR_DOCT sd 
WHERE kod = 1 AND dock = 1
AND bin_and(v_disp, ?) = ?
and kod_dock <> 888888
order by fio, im, ot`)
	INFO.Println(sql)
	stmt, err := m.DB.Prepare(sql)
	rows, err := stmt.Query(find.Name, find.Name)
	if err != nil {
		return nil, err
	}
	data := make([]types.Doctor, 0)
	for rows.Next() {
		var row types.Doctor
		err = rows.Scan(&row.Id, &row.Lname, &row.Fname, &row.Sname, &row.Access, &row.Z152)
		if err != nil {
			return nil, err
		}
		row.Lname, _ = utils.ToUTF8(strings.Trim(row.Lname, " "))
		row.Fname, _ = utils.ToUTF8(strings.Trim(row.Fname, " "))
		row.Sname, _ = utils.ToUTF8(strings.Trim(row.Sname, " "))
		data = append(data, row)
	}
	return &data, nil
}
