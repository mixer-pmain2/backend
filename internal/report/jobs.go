package report

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"pmain2/internal/consts"
	"pmain2/internal/models"
	"pmain2/pkg/cache"
	"pmain2/pkg/excel"
	"pmain2/pkg/utils"
	"strings"
	"time"
)

var (
	cacheReport = cache.CreateCache(time.Minute, time.Minute)
)

type ReportsJob struct {
}

// ReceptionLog Журнал приема
func (r *ReportsJob) ReceptionLog(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	var unit1, unit2 int
	if p.Unit == 1 {
		unit1 = 0
		unit2 = 1
	} else {
		unit1 = p.Unit
		unit2 = p.Unit + 1
	}
	sqlQuery := fmt.Sprintf(`select id_pat, f, b_day, diag, c1, mask, prichina, uch, uch_d
from f39_day('%s', %v) where (mask2 = %v) or (mask2 = %v)`, p.Filters.DateStart, p.UserId, unit1, unit2)
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	excelPage := excel.Page{f, sheet}

	m := models.Model.User
	doct, _ := m.Get(p.UserId, tx)
	dateStart, _ := time.Parse("2006-01-02", p.Filters.DateStart)

	sty := excelPage.CellStyleTitle("center", "middle", false, 9)
	excelPage.Title("Журнал посещения пациентов", cellExcel(1, 1), cellExcel(10, 1), sty)
	excelPage.Title(fmt.Sprintf("за %s   Врач - %s %s %s", dateStart.Format("02.01.2006"), doct.Lname, doct.Fname, doct.Sname), cellExcel(1, 2), cellExcel(10, 2), sty)

	f.SetCellStyle(sheet, cellExcel(1, 3), cellExcel(10, 3), excelPage.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 3), "#")
	f.SetCellStr(sheet, cellExcel(2, 3), "Шифр")
	f.SetCellStr(sheet, cellExcel(3, 3), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(4, 3), "Д.р.")
	f.SetCellStr(sheet, cellExcel(5, 3), "Уч.")
	f.SetCellStr(sheet, cellExcel(6, 3), "Диагноз")
	f.SetCellStr(sheet, cellExcel(7, 3), "Учет")
	f.SetCellStr(sheet, cellExcel(8, 3), "Где")
	f.SetCellStr(sheet, cellExcel(9, 3), "Уч.пр.")
	f.SetCellStr(sheet, cellExcel(10, 3), "Причина")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 28/7)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 50/7)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 120/7)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 66/7)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 36/7)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 57/7)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 37/7)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 43/7)
	f.SetColWidth(sheet, toCharStrConst(9), toCharStrConst(9), 35/7)
	f.SetColWidth(sheet, toCharStrConst(10), toCharStrConst(10), 112/7)

	nRow := 4
	nz := 1
	for rows.Next() {
		row := struct {
			PatientId   int
			PatientName string
			BDay        string
			Diagnose    string
			Category    string
			UnitName    string
			Reason      string
			Section     string
			SectionFrom string
		}{}
		err := rows.Scan(&row.PatientId, &row.PatientName, &row.BDay, &row.Diagnose, &row.Category, &row.UnitName, &row.Reason, &row.Section, &row.SectionFrom)
		if err != nil {
			ERROR.Println(err)
			return nil, err
		}
		bday, _ := time.Parse(time.RFC3339, row.BDay)
		diagnose, _ := utils.ToUTF8(row.Diagnose)
		category, _ := utils.ToUTF8(row.Category)
		unitName, _ := utils.ToUTF8(row.UnitName)
		reason, _ := utils.ToUTF8(row.Reason)
		patientName, _ := utils.ToUTF8(row.PatientName)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(10, nRow), excelPage.CellStyleBody(9))

		f.SetCellInt(sheet, cellExcel(1, nRow), nz)
		f.SetCellInt(sheet, cellExcel(2, nRow), row.PatientId)
		f.SetCellStr(sheet, cellExcel(3, nRow), patientName)
		f.SetCellStr(sheet, cellExcel(4, nRow), bday.Format("02.01.2006"))
		f.SetCellStr(sheet, cellExcel(5, nRow), strings.Trim(row.Section, " "))
		f.SetCellStr(sheet, cellExcel(6, nRow), diagnose)
		f.SetCellStr(sheet, cellExcel(7, nRow), category)
		f.SetCellStr(sheet, cellExcel(8, nRow), unitName)
		f.SetCellStr(sheet, cellExcel(9, nRow), strings.Trim(row.SectionFrom, " "))
		f.SetCellStr(sheet, cellExcel(10, nRow), reason)

		nz += 1
		nRow += 1
	}

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func (r *ReportsJob) VisitsPerPeriod(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeDate) < 2 {
		return nil, errors.New(consts.ArrErrors[750])
	}
	dateStart, _ := time.Parse("2006-01-02", p.Filters.RangeDate[0])
	dateEnd, _ := time.Parse("2006-01-02", p.Filters.RangeDate[1])

	sqlQuery := fmt.Sprintf(`select DAT, kol, prof from f39_diaposon('%s', '%s', %v, %v)`, p.Filters.RangeDate[0], p.Filters.RangeDate[1], p.UserId, p.Unit)
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	doct, _ := models.Model.User.Get(p.UserId, tx)
	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	page := excel.Page{f, sheet}

	sty := page.CellStyleTitle("center", "middle", false, 9)
	// row 1
	page.Title("Посещения врача", cellExcel(1, 1), cellExcel(8, 1), sty)
	// row 2
	page.Title(fmt.Sprintf("за период с %s по %s. Врач - %s %s %s",
		dateStart.Format("02.01.2006"), dateEnd.Format("02.01.2006"), doct.Lname, doct.Fname, doct.Sname),
		cellExcel(1, 2), cellExcel(8, 2), sty)
	// row 3
	f.SetCellStyle(sheet, cellExcel(1, 3), cellExcel(3, 3), page.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 3), "Дата")
	f.SetCellStr(sheet, cellExcel(2, 3), "Количество")
	f.SetCellStr(sheet, cellExcel(3, 3), "Профилактические")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 64/7)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 78/7)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 113/7)

	nRow := 4
	sumCount := 0
	sumProf := 0
	for rows.Next() {
		row := struct {
			date  string
			count int
			prof  int
		}{}
		rows.Scan(&row.date, &row.count, &row.prof)

		date, _ := time.Parse(time.RFC3339, row.date)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(3, nRow), page.CellStyleBody(9))

		page.SetCellStr(1, nRow, date.Format("02.01.2006"))
		page.SetCellInt(2, nRow, row.count)
		page.SetCellInt(3, nRow, row.prof)

		sumCount += row.count
		sumProf += row.prof
		nRow += 1
	}
	page.SetCellStr(1, nRow, "Всего")
	page.SetCellInt(2, nRow, sumCount)
	page.SetCellInt(3, nRow, sumProf)

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func (r *ReportsJob) AdmittedToTheHospital(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeSection) < 2 {
		return nil, errors.New(consts.ArrErrors[751])
	}
	if len(p.Filters.RangeDate) < 2 {
		return nil, errors.New(consts.ArrErrors[750])
	}
	u1 := p.Filters.RangeSection[0]
	u2 := p.Filters.RangeSection[1]
	if u2 < u1 {
		u2 = u1
	}

	sqlQuery := fmt.Sprintf(`select pi, fam, bd, otd, pdate, diag, diag_u, string, uch from SPS_INSTAC2(%v, %v,'%s', '%s')`,
		u1, u2, p.Filters.RangeDate[0], p.Filters.RangeDate[1])
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	page := excel.Page{f, sheet}

	dateStart, _ := time.Parse("2006-01-02", p.Filters.RangeDate[0])
	dateEnd, _ := time.Parse("2006-01-02", p.Filters.RangeDate[1])

	sty := page.CellStyleTitle("center", "middle", false, 9)
	// row 1
	page.Title("Список пациентов поступивших в стационар", cellExcel(1, 1), cellExcel(10, 1), sty)
	// row 2
	title2 := ""
	if u2 > u1 {
		title2 = fmt.Sprintf("уч. %v - %v", u1, u2)
	} else {
		title2 = fmt.Sprintf("уч. %v", u1)
	}
	page.Title(fmt.Sprintf("за период с %s по %s. %s",
		dateStart.Format("02.01.2006"), dateEnd.Format("02.01.2006"), title2),
		cellExcel(1, 2), cellExcel(10, 2), sty)
	// row 3
	f.SetCellStyle(sheet, cellExcel(1, 3), cellExcel(10, 3), page.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 3), "#")
	f.SetCellStr(sheet, cellExcel(2, 3), "Шифр")
	f.SetCellStr(sheet, cellExcel(3, 3), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(4, 3), "Д.р.")
	f.SetCellStr(sheet, cellExcel(5, 3), "Отд.")
	f.SetCellStr(sheet, cellExcel(6, 3), "Дата пост.")
	f.SetCellStr(sheet, cellExcel(7, 3), "Стационар")
	f.SetCellStr(sheet, cellExcel(8, 3), "Учет")
	f.SetCellStr(sheet, cellExcel(9, 3), "Причина")
	f.SetCellStr(sheet, cellExcel(10, 3), "Участок")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 21/7)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 53/7)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 127/7)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 67/7)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 42/7)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 67/7)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 68/7)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 68/7)
	f.SetColWidth(sheet, toCharStrConst(9), toCharStrConst(9), 65/7)
	f.SetColWidth(sheet, toCharStrConst(10), toCharStrConst(10), 42/7)

	nRow := 4
	nz := 1
	for rows.Next() {
		row := struct {
			patientId   int
			patientName string
			bday        string
			sectionStac int
			datePost    string
			diagStac    string
			diagReg     string
			reason      string
			sectionReg  int
		}{}
		rows.Scan(&row.patientId, &row.patientName, &row.bday, &row.sectionStac, &row.datePost, &row.diagStac, &row.diagReg, &row.reason, &row.sectionReg)

		bday, _ := time.Parse(time.RFC3339, row.bday)
		datePost, _ := time.Parse(time.RFC3339, row.datePost)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(10, nRow), page.CellStyleBody(9))
		patientName, _ := utils.ToUTF8(row.patientName)
		reason, _ := utils.ToUTF8(row.reason)
		page.SetCellInt(1, nRow, nz)
		page.SetCellInt(2, nRow, row.patientId)
		page.SetCellStr(3, nRow, patientName)
		page.SetCellStr(4, nRow, bday.Format("02.01.2006"))
		page.SetCellInt(5, nRow, row.sectionStac)
		page.SetCellStr(6, nRow, datePost.Format("02.01.2006"))
		page.SetCellStr(7, nRow, row.diagStac)
		page.SetCellStr(8, nRow, row.diagReg)
		page.SetCellStr(9, nRow, reason)
		page.SetCellInt(10, nRow, row.sectionReg)

		nRow += 1
		nz += 1
	}
	sty = page.CellStyleTitle("left", "middle", false, 9)
	page.Title("По данным отдела АСУ", cellExcel(1, nRow+1), cellExcel(5, nRow+1), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+2), cellExcel(5, nRow+2), sty)

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func (r *ReportsJob) DischargedFromHospital(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeDate) < 2 {
		return nil, errors.New(consts.ArrErrors[750])
	}
	if len(p.Filters.RangeSection) < 2 {
		return nil, errors.New(consts.ArrErrors[751])
	}
	u1 := p.Filters.RangeSection[0]
	u2 := p.Filters.RangeSection[1]
	if u2 < u1 {
		u2 = u1
	}

	sqlQuery := fmt.Sprintf(`select pi, otd, ni, fam, bd, pdate, exit_d, diag, diag_u, string from SPS_VSTAC ('%s', '%s', %v, %v) order by fam`,
		p.Filters.RangeDate[0], p.Filters.RangeDate[1], u1, u2)
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		patientId     int
		sectionStac   int
		historyNumber int
		patientName   string
		bday          string
		dateStart     string
		dateEnd       string
		diagStac      string
		diagReg       string
		reason        string
	}

	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	page := excel.Page{f, sheet}

	dateStart, _ := time.Parse("2006-01-02", p.Filters.RangeDate[0])
	dateEnd, _ := time.Parse("2006-01-02", p.Filters.RangeDate[1])

	sty := page.CellStyleTitle("center", "middle", false, 9)
	// row 1
	page.Title("Список выписанных из стационара", cellExcel(1, 1), cellExcel(10, 1), sty)
	// row 2
	title2 := ""
	if u2 > u1 {
		title2 = fmt.Sprintf("уч. %v - %v", u1, u2)
	} else {
		title2 = fmt.Sprintf("уч. %v", u1)
	}
	page.Title(fmt.Sprintf("за период с %s по %s. %s",
		dateStart.Format("02.01.2006"), dateEnd.Format("02.01.2006"), title2),
		cellExcel(1, 2), cellExcel(10, 2), sty)
	// row 3
	f.SetCellStyle(sheet, cellExcel(1, 3), cellExcel(10, 3), page.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 3), "Шифр")
	f.SetCellStr(sheet, cellExcel(2, 3), "№ Отд.")
	f.SetCellStr(sheet, cellExcel(3, 3), "№ Ист.")
	f.SetCellStr(sheet, cellExcel(4, 3), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(5, 3), "Д.р.")
	f.SetCellStr(sheet, cellExcel(6, 3), "Дата пост.")
	f.SetCellStr(sheet, cellExcel(7, 3), "Дата вып.")
	f.SetCellStr(sheet, cellExcel(8, 3), "Диагноз")
	f.SetCellStr(sheet, cellExcel(9, 3), "Наблюдение")
	f.SetCellStr(sheet, cellExcel(10, 3), "Причина")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 50/7)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 40/7)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 40/7)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 120/7)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 67/7)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 67/7)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 67/7)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 67/7)
	f.SetColWidth(sheet, toCharStrConst(9), toCharStrConst(9), 54/7)
	f.SetColWidth(sheet, toCharStrConst(10), toCharStrConst(10), 68/7)

	nRow := 4
	nz := 1
	for rows.Next() {
		row := data{}
		rows.Scan(&row.patientId, &row.sectionStac, &row.historyNumber, &row.patientName, &row.bday, &row.dateStart, &row.dateEnd, &row.diagStac, &row.diagReg, &row.reason)

		bday, _ := time.Parse(time.RFC3339, row.bday)
		dateStart, _ := time.Parse(time.RFC3339, row.dateStart)
		dateEnd, _ := time.Parse(time.RFC3339, row.dateEnd)
		patientName, _ := utils.ToUTF8(row.patientName)
		reason, _ := utils.ToUTF8(row.reason)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(10, nRow), page.CellStyleBody(9))
		page.SetCellInt(1, nRow, row.patientId)
		page.SetCellInt(2, nRow, row.sectionStac)
		page.SetCellInt(3, nRow, row.historyNumber)
		page.SetCellStr(4, nRow, patientName)
		page.SetCellStr(5, nRow, bday.Format("02.01.2006"))
		page.SetCellStr(6, nRow, dateStart.Format("02.01.2006"))
		page.SetCellStr(7, nRow, dateEnd.Format("02.01.2006"))
		page.SetCellStr(8, nRow, row.diagStac)
		page.SetCellStr(9, nRow, row.diagReg)
		page.SetCellStr(10, nRow, reason)

		nRow += 1
		nz += 1
	}
	sty = page.CellStyleTitle("left", "middle", false, 9)
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(5, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(5, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func (r *ReportsJob) Unvisited(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeSection) < 2 {
		return nil, errors.New(consts.ArrErrors[751])
	}

	u1 := p.Filters.RangeSection[0]
	u2 := p.Filters.RangeSection[1]
	if u2 < u1 {
		u2 = u1
	}
	if u1 == 0 {
		return nil, errors.New(consts.ArrErrors[752])
	}

	cacheName := fmt.Sprintf("Unvisited_%s_%v_%v", p.Filters.DateStart, u1, u2)
	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		return item.(*bytes.Buffer), nil
	}

	sqlQuery := fmt.Sprintf(`select pid, fio, bd, reg_date, categ, max_visit, diag, otd, exit_date, adr, max_v from spisok_nepos('%s', %v, %v)`,
		p.Filters.DateStart, u1, u2)
	if p.Filters.TypeCategory == "d" {
		sqlQuery += " where categ > 0 and categ < 10"
	}
	if p.Filters.TypeCategory == "k" {
		sqlQuery += " where categ = 10"
	}
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		patientId         int
		patientName       string
		bday              string
		dateReg           string
		category          string
		dateLastVisit     sql.NullString
		diag              string
		sectionStac       sql.NullString
		dateEnd           sql.NullString
		address           string
		dateLastVisitDocs sql.NullString
	}

	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	f.SetPageLayout(sheet,
		excelize.PageLayoutOrientation(excelize.OrientationLandscape),
		excelize.PageLayoutPaperSize(9),
	)
	page := excel.Page{f, sheet}

	dateStart, _ := time.Parse("2006-01-02", p.Filters.DateStart)

	sty := page.CellStyleTitle("center", "middle", false, 9)
	// row 1
	categoryType := ""
	if p.Filters.TypeCategory == "k" {
		categoryType = "К учет"
	}
	if p.Filters.TypeCategory == "d" {
		categoryType = "Д учет"
	}
	page.Title(fmt.Sprintf("Список непосещенных по уч. %v на %s %s", u1, dateStart.Format("02.01.2006"), categoryType), cellExcel(1, 1), cellExcel(11, 1), sty)
	// row 2
	f.SetCellStyle(sheet, cellExcel(1, 2), cellExcel(11, 2), page.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 2), "Шифр")
	f.SetCellStr(sheet, cellExcel(2, 2), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(3, 2), "Д.р.")
	f.SetCellStr(sheet, cellExcel(4, 2), "Дата взятия")
	f.SetCellStr(sheet, cellExcel(5, 2), "Группа")
	f.SetCellStr(sheet, cellExcel(6, 2), "Посл. посещ.")
	f.SetCellStr(sheet, cellExcel(7, 2), "Диагноз")
	f.SetCellStr(sheet, cellExcel(8, 2), "Отд.")
	f.SetCellStr(sheet, cellExcel(9, 2), "Выбыл")
	f.SetCellStr(sheet, cellExcel(10, 2), "Адрес")
	f.SetCellStr(sheet, cellExcel(11, 2), "Прочее")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 6.29)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 16.29)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 10.14)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 10.14)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 6.29)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 10.14)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 6.29)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 5.14)
	f.SetColWidth(sheet, toCharStrConst(9), toCharStrConst(9), 10.14)
	f.SetColWidth(sheet, toCharStrConst(10), toCharStrConst(10), 31.29)
	f.SetColWidth(sheet, toCharStrConst(11), toCharStrConst(11), 10.14)

	nRow := 3
	for rows.Next() {
		row := data{}
		rows.Scan(&row.patientId, &row.patientName, &row.bday, &row.dateReg, &row.category, &row.dateLastVisit, &row.diag, &row.sectionStac, &row.dateEnd, &row.address, &row.dateLastVisitDocs)

		bday, _ := time.Parse(time.RFC3339, row.bday)
		dateStart, _ := time.Parse(time.RFC3339, row.dateReg)

		dateEnd := row.dateEnd.String
		if dateEnd != "" {
			de, _ := time.Parse(time.RFC3339, dateEnd)
			dateEnd = de.Format("02.01.2006")
		}
		dateLastVisit := row.dateLastVisit.String
		if dateLastVisit != "" {
			dlv, _ := time.Parse(time.RFC3339, dateLastVisit)
			dateLastVisit = dlv.Format("02.01.2006")
		}
		dateLastVisitDocs := row.dateLastVisitDocs.String
		if dateLastVisitDocs != "" {
			dlv, _ := time.Parse(time.RFC3339, dateLastVisitDocs)
			dateLastVisitDocs = dlv.Format("02.01.2006")
		}
		patientName, _ := utils.ToUTF8(row.patientName)
		address, _ := utils.ToUTF8(row.address)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(11, nRow), page.CellStyleBody(9))
		page.SetCellInt(1, nRow, row.patientId)
		page.SetCellStr(2, nRow, patientName)
		page.SetCellStr(3, nRow, bday.Format("02.01.2006"))
		page.SetCellStr(4, nRow, dateStart.Format("02.01.2006"))
		page.SetCellStr(5, nRow, row.category)
		page.SetCellStr(6, nRow, dateLastVisit)
		page.SetCellStr(7, nRow, row.diag)
		page.SetCellStr(8, nRow, row.sectionStac.String)
		page.SetCellStr(9, nRow, dateEnd)
		page.SetCellStr(10, nRow, address)
		page.SetCellStr(11, nRow, dateLastVisitDocs)

		nRow += 1
	}
	sty = page.CellStyleTitle("left", "middle", false, 9)
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(5, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(5, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	cache.AppCache.Set(cacheName, buf, 0)
	return buf, nil
}

func (r *ReportsJob) Registered(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeDate) < 2 {
		return nil, errors.New(consts.ArrErrors[750])
	}
	if len(p.Filters.RangeSection) < 2 {
		return nil, errors.New(consts.ArrErrors[751])
	}

	u1 := p.Filters.RangeSection[0]
	u2 := p.Filters.RangeSection[1]
	if u2 < u1 {
		u2 = u1
	}
	if u1 == 0 {
		return nil, errors.New(consts.ArrErrors[752])
	}
	typeCategory := 0
	if p.Filters.TypeCategory == "k" {
		typeCategory = 1
	}

	cacheName := fmt.Sprintf("Registered_%v_%s_%s_%v_%v", u1, p.Filters.RangeDate[0], p.Filters.RangeDate[1], typeCategory, u2)
	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		return item.(*bytes.Buffer), nil
	}

	sqlQuery := fmt.Sprintf(`select pid, fio, bd, reg_date, categ, diag, adr, uch from spisok_in_Uchet(%v, '%s', '%s', %v, %v)`,
		u1, p.Filters.RangeDate[0], p.Filters.RangeDate[1], typeCategory, u2)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		patientId   sql.NullInt32
		patientName string
		bday        sql.NullString
		dateReg     sql.NullString
		category    sql.NullString
		diagnose    sql.NullString
		address     sql.NullString
		section     sql.NullInt32
	}

	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	f.SetPageLayout(sheet,
		excelize.PageLayoutOrientation(excelize.OrientationLandscape),
		excelize.PageLayoutPaperSize(9),
	)
	page := excel.Page{f, sheet}

	dateStart, _ := time.Parse("2006-01-02", p.Filters.RangeDate[0])
	dateEnd, _ := time.Parse("2006-01-02", p.Filters.RangeDate[1])

	sty := page.CellStyleTitle("center", "middle", false, 9)
	// row 1
	categoryType := ""
	if typeCategory == 1 {
		categoryType = "К"
	}
	if typeCategory == 0 {
		categoryType = "Д"
	}
	sectionString := ""
	if u1 < u2 {
		sectionString = fmt.Sprintf("%v - %v", u1, u2)
	} else {
		sectionString = fmt.Sprintf("%v", u1)
	}
	page.Title(
		fmt.Sprintf("Взятые под налблюдение группы \"%s\" с %s по %s на уч. %s",
			categoryType, dateStart.Format("02.01.2006"), dateEnd.Format("02.01.2006"), sectionString),
		cellExcel(1, 1), cellExcel(8, 1), sty)
	// row 2
	f.SetCellStyle(sheet, cellExcel(1, 2), cellExcel(8, 2), page.CellStyleHeader2(9))
	f.SetCellStr(sheet, cellExcel(1, 2), "")
	f.SetCellStr(sheet, cellExcel(2, 2), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(3, 2), "Д.р.")
	f.SetCellStr(sheet, cellExcel(4, 2), "Дата взятия")
	f.SetCellStr(sheet, cellExcel(5, 2), "Участок")
	f.SetCellStr(sheet, cellExcel(6, 2), "Группа")
	f.SetCellStr(sheet, cellExcel(7, 2), "Диагноз")
	f.SetCellStr(sheet, cellExcel(8, 2), "Адрес")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 6.29)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 28)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 10.14)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 10.14)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 7.29)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 8.14)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 9.29)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 42.14)

	nRow := 3
	for rows.Next() {
		row := data{}
		rows.Scan(&row.patientId, &row.patientName, &row.bday, &row.dateReg, &row.category, &row.diagnose, &row.address, &row.section)
		patientId := ""
		if row.patientId.Int32 != 0 {
			patientId = fmt.Sprintf("%v", row.patientId.Int32)
		}
		bday := ""
		if row.bday.String != "" {
			bd, _ := time.Parse(time.RFC3339, row.bday.String)
			bday = bd.Format("02.01.2006")
		}
		dateStart := ""
		if row.dateReg.String != "" {
			bd, _ := time.Parse(time.RFC3339, row.dateReg.String)
			dateStart = bd.Format("02.01.2006")
		}
		section := ""
		if row.section.Int32 != 0 {
			section = fmt.Sprintf("%v", row.section.Int32)
		}

		patientName, _ := utils.ToUTF8(row.patientName)
		address, _ := utils.ToUTF8(row.address.String)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(8, nRow), page.CellStyleBody2(9))
		page.SetCellStr(1, nRow, patientId)
		page.SetCellStr(2, nRow, patientName)
		page.SetCellStr(3, nRow, bday)
		page.SetCellStr(4, nRow, dateStart)
		page.SetCellStr(5, nRow, section)
		page.SetCellStr(6, nRow, row.category.String)
		page.SetCellStr(7, nRow, row.diagnose.String)
		page.SetCellStr(8, nRow, address)

		nRow += 1
	}
	sty = page.CellStyleTitle("left", "middle", false, 9)
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(8, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(8, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	cache.AppCache.Set(cacheName, buf, 0)
	return buf, nil
}

func (r *ReportsJob) Deregistered(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeDate) < 2 {
		return nil, errors.New(consts.ArrErrors[750])
	}
	if len(p.Filters.RangeSection) < 2 {
		return nil, errors.New(consts.ArrErrors[751])
	}

	u1 := p.Filters.RangeSection[0]
	u2 := p.Filters.RangeSection[1]
	if u2 < u1 {
		u2 = u1
	}
	if u1 == 0 {
		return nil, errors.New(consts.ArrErrors[752])
	}
	typeCategory := 0
	if p.Filters.TypeCategory == "k" {
		typeCategory = 1
	}

	cacheName := fmt.Sprintf("Deregistered_%s_%s_%v_%v_%v", p.Filters.RangeDate[0], p.Filters.RangeDate[1], u1, u2, typeCategory)
	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		return item.(*bytes.Buffer), nil
	}

	sqlQuery := fmt.Sprintf(`select l_patient_id, l_fio, l_bday, l_reg_date, l_categ_uch, l_diagnos, l_prich, l_uch from spisok_Out_uchet('%s', '%s', %v, %v, %v)`,
		p.Filters.RangeDate[0], p.Filters.RangeDate[1], u1, u2, typeCategory)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		patientId   sql.NullInt32
		patientName string
		bday        sql.NullString
		dateReg     sql.NullString
		category    sql.NullString
		diagnose    sql.NullString
		reason      sql.NullString
		section     sql.NullInt32
	}

	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	f.SetPageLayout(sheet,
		excelize.PageLayoutOrientation(excelize.OrientationPortrait),
		excelize.PageLayoutPaperSize(9),
	)
	page := excel.Page{f, sheet}

	dateStart, _ := time.Parse("2006-01-02", p.Filters.RangeDate[0])
	dateEnd, _ := time.Parse("2006-01-02", p.Filters.RangeDate[1])

	sty := page.CellStyleTitle("center", "middle", false, 9)
	// row 1
	categoryType := ""
	if typeCategory == 1 {
		categoryType = "К"
	}
	if typeCategory == 0 {
		categoryType = "Д"
	}
	sectionString := ""
	if u1 < u2 {
		sectionString = fmt.Sprintf("%v - %v", u1, u2)
	} else {
		sectionString = fmt.Sprintf("%v", u1)
	}
	page.Title(
		fmt.Sprintf("Снятые с \"%s\" наблюдения за период %s по %s с уч. %s",
			categoryType, dateStart.Format("02.01.2006"), dateEnd.Format("02.01.2006"), sectionString),
		cellExcel(1, 1), cellExcel(8, 1), sty)
	// row 2
	f.SetCellStyle(sheet, cellExcel(1, 2), cellExcel(8, 2), page.CellStyleHeader2(9))
	f.SetCellStr(sheet, cellExcel(1, 2), "")
	f.SetCellStr(sheet, cellExcel(2, 2), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(3, 2), "Д.р.")
	f.SetCellStr(sheet, cellExcel(4, 2), "Дата снятия")
	f.SetCellStr(sheet, cellExcel(5, 2), "Участок")
	f.SetCellStr(sheet, cellExcel(6, 2), "Группа")
	f.SetCellStr(sheet, cellExcel(7, 2), "Диагноз")
	f.SetCellStr(sheet, cellExcel(8, 2), "Причина")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 6.29)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 28)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 10.14)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 10.14)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 7.29)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 8.14)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 9.29)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 15.14)

	nRow := 3
	for rows.Next() {
		row := data{}
		rows.Scan(&row.patientId, &row.patientName, &row.bday, &row.dateReg, &row.category, &row.diagnose, &row.reason, &row.section)
		patientId := ""
		if row.patientId.Int32 != 0 {
			patientId = fmt.Sprintf("%v", row.patientId.Int32)
		}
		bday := ""
		if row.bday.String != "" {
			bd, _ := time.Parse(time.RFC3339, row.bday.String)
			bday = bd.Format("02.01.2006")
		}
		dateStart := ""
		if row.dateReg.String != "" {
			bd, _ := time.Parse(time.RFC3339, row.dateReg.String)
			dateStart = bd.Format("02.01.2006")
		}
		section := ""
		if row.section.Int32 != 0 {
			section = fmt.Sprintf("%v", row.section.Int32)
		}

		patientName, _ := utils.ToUTF8(row.patientName)
		reason, _ := utils.ToUTF8(row.reason.String)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(8, nRow), page.CellStyleBody2(9))
		page.SetCellStr(1, nRow, patientId)
		page.SetCellStr(2, nRow, patientName)
		page.SetCellStr(3, nRow, bday)
		page.SetCellStr(4, nRow, dateStart)
		page.SetCellStr(5, nRow, section)
		page.SetCellStr(6, nRow, row.category.String)
		page.SetCellStr(7, nRow, row.diagnose.String)
		page.SetCellStr(8, nRow, reason)

		nRow += 1
	}
	sty = page.CellStyleTitle("left", "middle", false, 9)
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(8, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(8, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	cache.AppCache.Set(cacheName, buf, 0)
	return buf, nil
}

func (r *ReportsJob) ConsistingOnTheSite(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeDate) < 2 {
		return nil, errors.New(consts.ArrErrors[750])
	}
	if len(p.Filters.RangeSection) < 2 {
		return nil, errors.New(consts.ArrErrors[751])
	}

	u1 := p.Filters.RangeSection[0]
	u2 := p.Filters.RangeSection[1]
	if u2 < u1 {
		u2 = u1
	}
	if u1 == 0 {
		return nil, errors.New(consts.ArrErrors[752])
	}
	typeCategory := 0
	if p.Filters.TypeCategory == "k" {
		typeCategory = 1
	}

	cacheName := fmt.Sprintf("ConsistingOnTheSite_%s_%s_%v_%v_%v", p.Filters.RangeDate[0], p.Filters.RangeDate[1], u1, u2, typeCategory)
	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		return item.(*bytes.Buffer), nil
	}

	sqlQuery := fmt.Sprintf(`select l_patient_id, l_fio, l_bday, l_reg_date, l_cat, l_diagnos, l_adres, l_primech, l_inv_beg from spisok_state_Uchet2(%v, '%s')`,
		u1, p.Filters.RangeDate[0])
	if p.Filters.TypeCategory == "k" {
		sqlQuery += " where l_cat = 10"
	}
	if p.Filters.TypeCategory == "d" {
		sqlQuery += " where l_cat > 0 and l_cat < 10"
	}
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		patientId    sql.NullInt32
		patientName  string
		bday         sql.NullString
		dateReg      sql.NullString
		category     sql.NullString
		diagnose     sql.NullString
		address      sql.NullString
		comment      sql.NullString
		dateInvStart sql.NullString
	}

	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	f.SetPageLayout(sheet,
		excelize.PageLayoutOrientation(excelize.OrientationLandscape),
		excelize.PageLayoutPaperSize(9),
	)
	page := excel.Page{f, sheet}

	dateStart, _ := time.Parse("2006-01-02", p.Filters.RangeDate[0])

	sty := page.CellStyleTitle("center", "middle", false, 9)
	// row 1
	categoryType := ""
	if typeCategory == 1 {
		categoryType = "на К учете"
	}
	if typeCategory == 0 {
		categoryType = "на Д учете"
	}
	page.Title(
		fmt.Sprintf("Список пациентов, находящихся по наблюдением на участке %v на %s %s",
			u1, dateStart.Format("02.01.2006"), categoryType),
		cellExcel(1, 1), cellExcel(9, 1), sty)
	// row 2
	f.SetCellStyle(sheet, cellExcel(1, 2), cellExcel(9, 2), page.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 2), "")
	f.SetCellStr(sheet, cellExcel(2, 2), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(3, 2), "Д.р.")
	f.SetCellStr(sheet, cellExcel(4, 2), "Дата взятия")
	f.SetCellStr(sheet, cellExcel(5, 2), "Группа")
	f.SetCellStr(sheet, cellExcel(6, 2), "Диагноз")
	f.SetCellStr(sheet, cellExcel(7, 2), "Адрес")
	f.SetCellStr(sheet, cellExcel(8, 2), "Инвалид.")
	f.SetCellStr(sheet, cellExcel(9, 2), "Примечание")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 6.29)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 29.57)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 10.14)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 10.14)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 7.29)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 8.14)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 37)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 10.14)
	f.SetColWidth(sheet, toCharStrConst(9), toCharStrConst(9), 10.14)

	nRow := 3
	sum := 0
	for rows.Next() {
		row := data{}
		rows.Scan(&row.patientId, &row.patientName, &row.bday, &row.dateReg, &row.category, &row.diagnose, &row.address, &row.comment, &row.dateInvStart)
		patientId := ""
		if row.patientId.Int32 != 0 {
			patientId = fmt.Sprintf("%v", row.patientId.Int32)
		}
		bday := ""
		if row.bday.String != "" {
			bd, _ := time.Parse(time.RFC3339, row.bday.String)
			bday = bd.Format("02.01.2006")
		}
		dateStart := ""
		if row.dateReg.String != "" {
			d, _ := time.Parse(time.RFC3339, row.dateReg.String)
			dateStart = d.Format("02.01.2006")
		}
		dateInvStart := ""
		if row.dateInvStart.String != "" {
			d, _ := time.Parse(time.RFC3339, row.dateInvStart.String)
			dateInvStart = d.Format("02.01.2006")
		}

		patientName, _ := utils.ToUTF8(row.patientName)
		address, _ := utils.ToUTF8(row.address.String)
		comment, _ := utils.ToUTF8(row.comment.String)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(9, nRow), page.CellStyleBody(9))
		page.SetCellStr(1, nRow, patientId)
		page.SetCellStr(2, nRow, patientName)
		page.SetCellStr(3, nRow, bday)
		page.SetCellStr(4, nRow, dateStart)
		page.SetCellStr(5, nRow, row.category.String)
		page.SetCellStr(6, nRow, row.diagnose.String)
		page.SetCellStr(7, nRow, address)
		page.SetCellStr(8, nRow, dateInvStart)
		page.SetCellStr(9, nRow, comment)

		nRow += 1
		sum += 1
	}
	sty = page.CellStyleTitle("left", "middle", false, 9)
	page.Title(fmt.Sprintf("Всего: %v", sum), cellExcel(1, nRow), cellExcel(9, nRow), sty)
	nRow += 1
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(9, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(8, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	cache.AppCache.Set(cacheName, buf, 0)
	return buf, nil
}
