package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pmain2/internal/apperror"
)

type SprDoctModel struct {
	Db *sql.DB
}

func (m *SprDoctModel) FoundByFIO(lname, fname, sname string) (*[]SprDoct, error) {
	data := []SprDoct{}
	sql := fmt.Sprintf(
		"select kod_dock_i, fio, im, ot FROM SPR_DOCT where position('%s', fio)>0 and position('%s', im)>0 and position('%s', ot)>0",
		lname, fname, sname)
	INFO.Println(sql)
	rows, err := m.Db.Query(sql)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	for rows.Next() {
		row := SprDoct{}
		err = rows.Scan(&row.Id, &row.Lname, &row.Fname, &row.Sname)
		if err != nil {
			ERROR.Println(err.Error())
			return nil, err
		}
		data = append(data, row)

	}
	res, _ := json.Marshal(&data)
	INFO.Println(string(res))
	return &data, nil
}

func (m *SprDoctModel) Get(id int) (*SprDoct, error) {
	data := SprDoct{}
	sql := fmt.Sprintf(
		"select kod_dock_i, fio, im, ot FROM SPR_DOCT where kod_dock_i=%v", id)
	INFO.Println(sql)
	row := m.Db.QueryRow(sql)
	err := row.Scan(&data.Id, &data.Lname, &data.Fname, &data.Sname)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}

	err = data.Serialize()
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}

	if data.Id == 0 {
		err = apperror.ErrDataNotFound
		return nil, err
	}

	res, err := json.Marshal(&data)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	INFO.Println(res)

	return &data, nil
}

func (m *SprDoctModel) UserAuth(login, password string) (bool, error) {
	var n int
	sql := fmt.Sprintf(
		"select count(*) FROM SPR_DOCT where kod_dock_i=%v and pass_new='%s'",
		login, password)
	INFO.Println(sql)
	rows := m.Db.QueryRow(sql)
	err := rows.Scan(&n)
	if err != nil {
		ERROR.Println(err.Error())
		return false, err
	}
	res, _ := json.Marshal(n)
	INFO.Println(string(res))
	return n > 0, nil
}

func (m *SprDoctModel) GetPrava(id int) (*map[int]int, error) {
	data := make(map[int]int, 0)
	sql := fmt.Sprintf("SELECT PODR, SUM(PRAVA)  FROM DOCK_PRAVA dp WHERE dp.KOD_DOCT=%v group by podr", id)
	INFO.Println(sql)
	rows, err := m.Db.Query(sql)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	for rows.Next() {
		var p, s int
		err = rows.Scan(&p, &s)
		if err != nil {
			return nil, err
		}
		data[p] = s
	}

	return &data, nil
}
