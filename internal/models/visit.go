package models

import (
	"database/sql"
	"fmt"
)

type VisitModel struct {
	DB *sql.DB
}

func createVisit() *VisitModel {
	return &VisitModel{}
}

func (m *VisitModel) GetVisits(patientId, numPage, pageSize int) (*[]Visit, error) {
	skip := numPage * pageSize
	sql := fmt.Sprintf(`SELECT FIRST %s skip %s 
V_NUM, PATIENT_ID, V_DATE, 
NAME_DOCT, DIAGNOSE, MASKA1, 
CASE WHEN MASKA2 in (0, 1) THEN 1 ELSE MASKA2 END podr, bin_and(1, MASKA2) home
FROM visit v 
WHERE PATIENT_ID = %s`, pageSize, skip, patientId)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	var data []Visit
	for rows.Next() {
		v := Visit{}
		err = rows.Scan(&v.Id, &v.PatientId,
			&v.Date, &v.DockId,
			&v.Diagnose, &v.Type,
			&v.Pord, &v.Home)
		if err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return &data, nil

}
