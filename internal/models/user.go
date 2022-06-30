package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pmain2/internal/apperror"
	"pmain2/internal/types"
)

type userModel struct {
	DB *sql.DB
}

func createUser(db *sql.DB) *userModel {
	return &userModel{DB: db}
}

func (m *userModel) FoundByFIO(lname, fname, sname string) (*[]SprDoct, error) {
	data := []SprDoct{}
	sql := fmt.Sprintf(
		"select kod_dock_i, fio, im, ot FROM SPR_DOCT where position('%s', fio)>0 and position('%s', im)>0 and position('%s', ot)>0",
		lname, fname, sname)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
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

func (m *userModel) Get(id int) (*SprDoct, error) {
	data := SprDoct{}
	sql := fmt.Sprintf(
		"select kod_dock_i, fio, im, ot FROM SPR_DOCT where kod_dock_i=%v", id)
	INFO.Println(sql)
	row := m.DB.QueryRow(sql)
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
		return nil, err
	}
	INFO.Println(string(res))

	return &data, nil
}

func (m *userModel) UserAuth(login, password string) (bool, error) {
	var n int
	sql := fmt.Sprintf(
		"select count(*) FROM SPR_DOCT where kod_dock_i=? and pass_new=?")
	INFO.Println(sql)
	stmt, err := m.DB.Prepare(sql)
	if err != nil {
		return false, err
	}
	rows := stmt.QueryRow(login, password)
	//rows := m.DB.QueryRow(sql)
	err = rows.Scan(&n)
	if err != nil {
		return false, err
	}
	res, _ := json.Marshal(n)
	INFO.Println(string(res))
	return n > 0, nil
}

func (m *userModel) GetPrava(id int) (*map[int]int, error) {
	data := make(map[int]int, 0)
	sql := fmt.Sprintf(`SELECT PODR, SUM(PRAVA) 
 FROM DOCK_PRAVA dp WHERE dp.KOD_DOCT=%v
 and 'NOW' BETWEEN DATE_N AND DATE_E 
 group by podr`, id)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
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

func (m *userModel) GetUch(id int) (*map[int][]int, error) {
	data := make(map[int][]int, 0)
	sql := fmt.Sprintf(`select PODRAZ, UCH from dock_uch
where dock = %v
and priz = 1
order by podraz, uch;`, id)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
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
		data[p] = append(data[p], s)
	}

	return &data, nil
}

func (m *userModel) ChangePassword(data types.ChangePassword) (sql.Result, error) {
	sql := fmt.Sprintf(`update spr_doct set pass_new ='%s' where kod_dock_i = %v`, data.NewPassword, data.UserId)
	INFO.Println(sql)
	return m.DB.Exec(sql)
}
