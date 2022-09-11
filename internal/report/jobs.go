package report

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"pmain2/internal/consts"
	"pmain2/internal/models"
	"pmain2/internal/types"
	"pmain2/pkg/cache"
	"pmain2/pkg/excel"
	"pmain2/pkg/utils"
	"strconv"
	"strings"
	"time"
)

const (
	HEIGHT_ROW  = 13
	HEIGHT_2ROW = 25
	HEIGHT_3ROW = 34
	HEIGHT_4ROW = 47
	HEIGHT_5ROW = 58
	HEIGHT_6ROW = 69
)

var (
	cacheReport = cache.CreateCache(time.Minute, time.Minute)
)

type ReportsJob struct {
}

// ReceptionLog Журнал приема
func (r *ReportsJob) ReceptionLog(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	var unit1, unit2 int
	if p.Filters.Unit == 1 {
		unit1 = 0
		unit2 = 1
	} else {
		unit1 = p.Filters.Unit
		unit2 = p.Filters.Unit + 1
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

	sty := excelPage.CellStyleTitle("center", "center", false, 9)
	excelPage.Title("Журнал посещения пациентов", cellExcel(1, 1), cellExcel(10, 1), sty)
	excelPage.Title(fmt.Sprintf("за %s   Врач - %s %s %s", dateStart.Format("02.01.2006"), doct.Lname, doct.Fname, doct.Sname), cellExcel(1, 2), cellExcel(10, 2), sty)
	f.SetRowHeight(sheet, 3, 27)
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
		excelPage.SetCellInt(2, nRow, row.PatientId)
		excelPage.SetCellStr(3, nRow, patientName)
		excelPage.SetCellStr(4, nRow, bday.Format("02.01.2006"))
		excelPage.SetCellStr(5, nRow, strings.Trim(row.Section, " "))
		excelPage.SetCellStr(6, nRow, diagnose)
		excelPage.SetCellStr(7, nRow, category)
		excelPage.SetCellStr(8, nRow, unitName)
		excelPage.SetCellStr(9, nRow, strings.Trim(row.SectionFrom, " "))
		excelPage.SetCellStr(10, nRow, reason)

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

	sqlQuery := fmt.Sprintf(`select DAT, kol, prof from f39_diaposon('%s', '%s', %v, %v)`, p.Filters.RangeDate[0], p.Filters.RangeDate[1], p.UserId, p.Filters.Unit)
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

	sty := page.CellStyleTitle("center", "center", false, 9)
	// row 1
	page.Title("Посещения врача", cellExcel(1, 1), cellExcel(8, 1), sty)
	// row 2
	page.Title(fmt.Sprintf("за период с %s по %s. Врач - %s %s %s",
		dateStart.Format("02.01.2006"), dateEnd.Format("02.01.2006"), doct.Lname, doct.Fname, doct.Sname),
		cellExcel(1, 2), cellExcel(8, 2), sty)
	// row 3
	f.SetRowHeight(sheet, 3, 27)
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

	sty := page.CellStyleTitle("center", "center", false, 9)
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
	f.SetRowHeight(sheet, 3, 27)
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
	sty = page.CellStyleTitle("left", "center", false, 9)
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

	sty := page.CellStyleTitle("center", "center", false, 9)
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
	f.SetRowHeight(sheet, 3, 27)
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
	sty = page.CellStyleTitle("left", "center", false, 9)
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

	sty := page.CellStyleTitle("center", "center", false, 9)
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
	f.SetRowHeight(sheet, 2, 27)
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
	sty = page.CellStyleTitle("left", "center", false, 9)
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

	if p.Filters.TypeCategory == "" {
		p.Filters.TypeCategory = "d"
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

	sty := page.CellStyleTitle("center", "center", false, 9)
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
	f.SetRowHeight(sheet, 2, 27)
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
	sty = page.CellStyleTitle("left", "center", false, 9)
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

	if p.Filters.TypeCategory == "" {
		p.Filters.TypeCategory = "d"
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

	sty := page.CellStyleTitle("center", "center", false, 9)
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
	f.SetRowHeight(sheet, 2, 27)
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
	sty = page.CellStyleTitle("left", "center", false, 9)
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

	sty := page.CellStyleTitle("center", "center", false, 9)
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
	f.SetRowHeight(sheet, 2, 27)
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
	sty = page.CellStyleTitle("left", "center", false, 9)
	page.Title(fmt.Sprintf("Всего: %v", sum), cellExcel(1, nRow), cellExcel(9, nRow), sty)
	nRow += 1
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(9, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(8, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	cache.AppCache.Set(cacheName, buf, 0)
	return buf, nil
}

func (r *ReportsJob) ThoseInTheHospital(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeSection) < 2 {
		return nil, errors.New(consts.ArrErrors[751])
	}
	//if len(p.Filters.RangeDate) < 2 {
	//	return nil, errors.New(consts.ArrErrors[750])
	//}
	u1 := p.Filters.RangeSection[0]
	u2 := p.Filters.RangeSection[1]
	if u2 < u1 {
		u2 = u1
	}

	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	page := excel.Page{f, sheet}

	sty := page.CellStyleTitle("center", "center", false, 9)
	// row 1
	page.Title("Список пациентов находящихся в стационар", cellExcel(1, 1), cellExcel(10, 1), sty)
	// row 2
	title2 := ""
	if u2 > u1 {
		title2 = fmt.Sprintf("уч. %v - %v", u1, u2)
	} else {
		title2 = fmt.Sprintf("уч. %v", u1)
	}
	page.Title(fmt.Sprintf("%s", title2),
		cellExcel(1, 2), cellExcel(10, 2), sty)
	// row 3
	f.SetRowHeight(sheet, 3, 27)
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

	sqlQuery := fmt.Sprintf(`SELECT pi, fam, bd, otd, pdate, diag, diag_u, string, uch  
from SPS_INSTAC(%v, %v) order by fam`, u1, u2)
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

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
	sty = page.CellStyleTitle("left", "center", false, 9)
	page.Title("По данным отдела АСУ", cellExcel(1, nRow+1), cellExcel(5, nRow+1), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+2), cellExcel(5, nRow+2), sty)

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func (r *ReportsJob) HospitalTreatment(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {

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
	cols := 14
	sty := page.CellStyleTitle("center", "center", false, 9)
	page.Title(
		fmt.Sprintf("Список находящихся на ПЛ по состоянию на %s", dateStart.Format("02.01.2006")),
		cellExcel(1, 1), cellExcel(cols, 1), sty)
	// row 2
	f.SetRowHeight(sheet, 2, 27)
	f.SetCellStyle(sheet, cellExcel(1, 2), cellExcel(cols, 2), page.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 2), "")
	f.SetCellStr(sheet, cellExcel(2, 2), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(3, 2), "Д.р.")
	f.SetCellStr(sheet, cellExcel(4, 2), "Первое определ.")
	f.SetCellStr(sheet, cellExcel(5, 2), "Последн определ.")
	f.SetCellStr(sheet, cellExcel(6, 2), "Мех")
	f.SetCellStr(sheet, cellExcel(7, 2), "Отд.")
	f.SetCellStr(sheet, cellExcel(8, 2), "Кат. учета")
	f.SetCellStr(sheet, cellExcel(9, 2), "Участок учета")
	f.SetCellStr(sheet, cellExcel(10, 2), "Группа учета")
	f.SetCellStr(sheet, cellExcel(11, 2), "Последний осмотр")
	f.SetCellStr(sheet, cellExcel(12, 2), "Адрес")
	f.SetCellStr(sheet, cellExcel(13, 2), "Инв.")
	f.SetCellStr(sheet, cellExcel(14, 2), "Заболел")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 7.29)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 18.57)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 9.43)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 9.43)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 9.43)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 7.43)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 6)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 6)
	f.SetColWidth(sheet, toCharStrConst(9), toCharStrConst(9), 6)
	f.SetColWidth(sheet, toCharStrConst(10), toCharStrConst(10), 6)
	f.SetColWidth(sheet, toCharStrConst(11), toCharStrConst(11), 9.43)
	f.SetColWidth(sheet, toCharStrConst(12), toCharStrConst(12), 13.43)
	f.SetColWidth(sheet, toCharStrConst(13), toCharStrConst(13), 8.57)
	f.SetColWidth(sheet, toCharStrConst(14), toCharStrConst(14), 8.57)

	sqlQuery := fmt.Sprintf(`SELECT patient_id, fio, bday, nach, tek_opr, 
CASE WHEN tr = 1 THEN 'п/п' when tr = 2 THEN 'н/л' ELSE '?' end, 
otd, kat, uch, gruppa, osmotr, adr, inv, 
CASE WHEN zabolel = 1 THEN 'до' WHEN zabolel = 2 THEN 'пос' ELSE '?' end 
from sps_pl_nahod('%s') where (gruppa <> 2 or (gruppa = 2 and otd is not null)) and gruppa <> 4 order by otd, fio`,
		p.Filters.DateStart)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		patientId   int
		patientName string
		bday        string
		startOp     string
		endOp       string
		mech        string
		sectionStac sql.NullInt64
		category    sql.NullInt64
		sectionReg  sql.NullInt64
		group       sql.NullInt64
		dateView    string
		address     string
		invalid     string
		sick        string
	}

	nRow := 3
	sum := 0
	for rows.Next() {
		rereg := false
		row := data{}
		rows.Scan(&row.patientId, &row.patientName, &row.bday, &row.startOp, &row.endOp, &row.mech, &row.sectionStac,
			&row.category, &row.sectionReg, &row.group, &row.dateView, &row.address, &row.invalid, &row.sick)
		bd, _ := time.Parse(time.RFC3339, row.bday)
		bday := bd.Format("02.01.2006")
		startOp := ""
		if row.startOp != "" {
			d, _ := time.Parse(time.RFC3339, row.startOp)
			startOp = d.Format("02.01.2006")
			if d.Sub(time.Now().AddDate(0, 0, -152)) < 0 {
				rereg = true
			}
		}
		endOp := ""
		if row.endOp != "" {
			d, _ := time.Parse(time.RFC3339, row.endOp)
			endOp = d.Format("02.01.2006")
		}
		sectionStac := ""
		if row.sectionStac.Int64 > 0 {
			sectionStac = strconv.FormatInt(row.sectionStac.Int64, 10)
		}
		category := ""
		if row.category.Int64 > 0 {
			category = strconv.FormatInt(row.category.Int64, 10)
		}
		sectionReg := ""
		if row.sectionReg.Int64 > 0 {
			sectionReg = strconv.FormatInt(row.sectionReg.Int64, 10)
		}
		group := ""
		if row.group.Int64 > 0 {
			group = strconv.FormatInt(row.group.Int64, 10)
		}
		dateView := ""
		if row.dateView != "" {
			d, _ := time.Parse(time.RFC3339, row.dateView)
			dateView = d.Format("02.01.2006")
			if d.Sub(time.Now().AddDate(0, 0, -152)) < 0 {
				rereg = true
			}
		}

		patientName, _ := utils.ToUTF8(row.patientName)
		address, _ := utils.ToUTF8(row.address)
		//mech, _ := utils.ToUTF8(row.mech)
		invalid, _ := utils.ToUTF8(row.invalid)
		//sick, _ := utils.ToUTF8(row.sick)
		if rereg {
			f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(cols, nRow), page.CellStyleBodyColor(9, "#ededed"))
		} else {
			f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(cols, nRow), page.CellStyleBody(9))
		}
		page.SetCellInt(1, nRow, row.patientId)
		page.SetCellStr(2, nRow, patientName)
		page.SetCellStr(3, nRow, bday)
		page.SetCellStr(4, nRow, startOp)
		page.SetCellStr(5, nRow, endOp)
		page.SetCellStr(6, nRow, row.mech)
		page.SetCellStr(7, nRow, sectionStac)
		page.SetCellStr(8, nRow, category)
		page.SetCellStr(9, nRow, sectionReg)
		page.SetCellStr(10, nRow, group)
		page.SetCellStr(11, nRow, dateView)
		page.SetCellStr(12, nRow, address)
		page.SetCellStr(13, nRow, invalid)
		page.SetCellStr(14, nRow, row.sick)

		nRow += 1
		sum += 1
	}
	sty = page.CellStyleTitle("left", "center", false, 9)
	page.Title(fmt.Sprintf("Всего: %v", sum), cellExcel(1, nRow), cellExcel(9, nRow), sty)
	nRow += 1
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(9, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(8, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func (r *ReportsJob) AmbulatoryTreatment(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {

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
	cols := 14
	sty := page.CellStyleTitle("center", "center", false, 9)
	page.Title(
		fmt.Sprintf("Список находящихся на АПЛ по состоянию на %s", dateStart.Format("02.01.2006")),
		cellExcel(1, 1), cellExcel(cols, 1), sty)
	// row 2
	f.SetRowHeight(sheet, 2, 27)
	f.SetCellStyle(sheet, cellExcel(1, 2), cellExcel(cols, 2), page.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 2), "")
	f.SetCellStr(sheet, cellExcel(2, 2), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(3, 2), "Д.р.")
	f.SetCellStr(sheet, cellExcel(4, 2), "Первое определ.")
	f.SetCellStr(sheet, cellExcel(5, 2), "Последн определ.")
	f.SetCellStr(sheet, cellExcel(6, 2), "Мех")
	f.SetCellStr(sheet, cellExcel(7, 2), "Отд.")
	f.SetCellStr(sheet, cellExcel(8, 2), "Кат. учета")
	f.SetCellStr(sheet, cellExcel(9, 2), "Участок учета")
	f.SetCellStr(sheet, cellExcel(10, 2), "Группа учета")
	f.SetCellStr(sheet, cellExcel(11, 2), "Последний осмотр")
	f.SetCellStr(sheet, cellExcel(12, 2), "Адрес")
	f.SetCellStr(sheet, cellExcel(13, 2), "Инв.")
	f.SetCellStr(sheet, cellExcel(14, 2), "Заболел")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 7.29)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 18.57)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 9.43)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 9.43)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 9.43)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 7.43)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 6)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 6)
	f.SetColWidth(sheet, toCharStrConst(9), toCharStrConst(9), 6)
	f.SetColWidth(sheet, toCharStrConst(10), toCharStrConst(10), 6)
	f.SetColWidth(sheet, toCharStrConst(11), toCharStrConst(11), 9.43)
	f.SetColWidth(sheet, toCharStrConst(12), toCharStrConst(12), 13.43)
	f.SetColWidth(sheet, toCharStrConst(13), toCharStrConst(13), 8.57)
	f.SetColWidth(sheet, toCharStrConst(14), toCharStrConst(14), 8.57)

	sqlQuery := fmt.Sprintf(`SELECT patient_id, fio, bday, nach, tek_opr, 
CASE WHEN tr = 1 THEN 'п/п' when tr = 2 THEN 'н/л' ELSE '?' end, 
otd, kat, uch, gruppa, osmotr, adr, inv, 
CASE WHEN zabolel = 1 THEN 'до' WHEN zabolel = 2 THEN 'пос' ELSE '?' end 
from sps_pl_nahod('%s') where (gruppa <> 2 or (gruppa = 2 and otd is not null)) and gruppa = 4 order by otd, fio`,
		p.Filters.DateStart)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		patientId   int
		patientName string
		bday        string
		startOp     string
		endOp       string
		mech        string
		sectionStac sql.NullInt64
		category    sql.NullInt64
		sectionReg  sql.NullInt64
		group       sql.NullInt64
		dateView    string
		address     string
		invalid     string
		sick        string
	}

	nRow := 3
	sum := 0
	for rows.Next() {
		rereg := false
		row := data{}
		rows.Scan(&row.patientId, &row.patientName, &row.bday, &row.startOp, &row.endOp, &row.mech, &row.sectionStac,
			&row.category, &row.sectionReg, &row.group, &row.dateView, &row.address, &row.invalid, &row.sick)
		bd, _ := time.Parse(time.RFC3339, row.bday)
		bday := bd.Format("02.01.2006")
		startOp := ""
		if row.startOp != "" {
			d, _ := time.Parse(time.RFC3339, row.startOp)
			startOp = d.Format("02.01.2006")
			if d.Sub(time.Now().AddDate(0, 0, -152)) < 0 {
				rereg = true
			}
		}
		endOp := ""
		if row.endOp != "" {
			d, _ := time.Parse(time.RFC3339, row.endOp)
			endOp = d.Format("02.01.2006")
		}
		sectionStac := ""
		if row.sectionStac.Int64 > 0 {
			sectionStac = strconv.FormatInt(row.sectionStac.Int64, 10)
		}
		category := ""
		if row.category.Int64 > 0 {
			category = strconv.FormatInt(row.category.Int64, 10)
		}
		sectionReg := ""
		if row.sectionReg.Int64 > 0 {
			sectionReg = strconv.FormatInt(row.sectionReg.Int64, 10)
		}
		group := ""
		if row.group.Int64 > 0 {
			group = strconv.FormatInt(row.group.Int64, 10)
		}
		dateView := ""
		if row.dateView != "" {
			d, _ := time.Parse(time.RFC3339, row.dateView)
			dateView = d.Format("02.01.2006")
			if d.Sub(time.Now().AddDate(0, 0, -152)) < 0 {
				rereg = true
			}
		}

		patientName, _ := utils.ToUTF8(row.patientName)
		address, _ := utils.ToUTF8(row.address)
		//mech, _ := utils.ToUTF8(row.mech)
		invalid, _ := utils.ToUTF8(row.invalid)
		//sick, _ := utils.ToUTF8(row.sick)

		if rereg {
			f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(cols, nRow), page.CellStyleBodyColor(9, "#ededed"))
		} else {
			f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(cols, nRow), page.CellStyleBody(9))
		}
		page.SetCellInt(1, nRow, row.patientId)
		page.SetCellStr(2, nRow, patientName)
		page.SetCellStr(3, nRow, bday)
		page.SetCellStr(4, nRow, startOp)
		page.SetCellStr(5, nRow, endOp)
		page.SetCellStr(6, nRow, row.mech)
		page.SetCellStr(7, nRow, sectionStac)
		page.SetCellStr(8, nRow, category)
		page.SetCellStr(9, nRow, sectionReg)
		page.SetCellStr(10, nRow, group)
		page.SetCellStr(11, nRow, dateView)
		page.SetCellStr(12, nRow, address)
		page.SetCellStr(13, nRow, invalid)
		page.SetCellStr(14, nRow, row.sick)

		nRow += 1
		sum += 1
	}
	sty = page.CellStyleTitle("left", "center", false, 9)
	page.Title(fmt.Sprintf("Всего: %v", sum), cellExcel(1, nRow), cellExcel(9, nRow), sty)
	nRow += 1
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(9, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(8, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func (r *ReportsJob) PBSTIN(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {

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
	cols := 14
	sty := page.CellStyleTitle("center", "center", false, 9)
	page.Title(
		fmt.Sprintf("Список находящихся в ПБСТИН по состоянию на %s", dateStart.Format("02.01.2006")),
		cellExcel(1, 1), cellExcel(cols, 1), sty)
	// row 2
	f.SetRowHeight(sheet, 2, 27)
	f.SetCellStyle(sheet, cellExcel(1, 2), cellExcel(cols, 2), page.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 2), "")
	f.SetCellStr(sheet, cellExcel(2, 2), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(3, 2), "Д.р.")
	f.SetCellStr(sheet, cellExcel(4, 2), "Первое определ.")
	f.SetCellStr(sheet, cellExcel(5, 2), "Последн определ.")
	f.SetCellStr(sheet, cellExcel(6, 2), "Мех")
	f.SetCellStr(sheet, cellExcel(7, 2), "Отд.")
	f.SetCellStr(sheet, cellExcel(8, 2), "Кат. учета")
	f.SetCellStr(sheet, cellExcel(9, 2), "Участок учета")
	f.SetCellStr(sheet, cellExcel(10, 2), "Группа учета")
	f.SetCellStr(sheet, cellExcel(11, 2), "Последний осмотр")
	f.SetCellStr(sheet, cellExcel(12, 2), "Адрес")
	f.SetCellStr(sheet, cellExcel(13, 2), "Инв.")
	f.SetCellStr(sheet, cellExcel(14, 2), "Заболел")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 7.29)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 18.57)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 9.43)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 9.43)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 9.43)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 7.43)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 6)
	f.SetColWidth(sheet, toCharStrConst(8), toCharStrConst(8), 6)
	f.SetColWidth(sheet, toCharStrConst(9), toCharStrConst(9), 6)
	f.SetColWidth(sheet, toCharStrConst(10), toCharStrConst(10), 6)
	f.SetColWidth(sheet, toCharStrConst(11), toCharStrConst(11), 9.43)
	f.SetColWidth(sheet, toCharStrConst(12), toCharStrConst(12), 13.43)
	f.SetColWidth(sheet, toCharStrConst(13), toCharStrConst(13), 8.57)
	f.SetColWidth(sheet, toCharStrConst(14), toCharStrConst(14), 8.57)

	sqlQuery := fmt.Sprintf(`SELECT patient_id, fio, bday, nach, tek_opr, 
CASE WHEN tr = 1 THEN 'п/п' when tr = 2 THEN 'н/л' ELSE '?' end, 
otd, kat, uch, gruppa, osmotr, adr, inv, 
CASE WHEN zabolel = 1 THEN 'до' WHEN zabolel = 2 THEN 'пос' ELSE '?' end 
from sps_pl_nahod('%s') where gruppa = 2 order by fio`,
		p.Filters.DateStart)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		patientId   int
		patientName string
		bday        string
		startOp     string
		endOp       string
		mech        string
		sectionStac sql.NullInt64
		category    sql.NullInt64
		sectionReg  sql.NullInt64
		group       sql.NullInt64
		dateView    string
		address     string
		invalid     string
		sick        string
	}

	nRow := 3
	sum := 0
	for rows.Next() {
		row := data{}
		rows.Scan(&row.patientId, &row.patientName, &row.bday, &row.startOp, &row.endOp, &row.mech, &row.sectionStac,
			&row.category, &row.sectionReg, &row.group, &row.dateView, &row.address, &row.invalid, &row.sick)
		bd, _ := time.Parse(time.RFC3339, row.bday)
		bday := bd.Format("02.01.2006")
		startOp := ""
		if row.startOp != "" {
			d, _ := time.Parse(time.RFC3339, row.startOp)
			startOp = d.Format("02.01.2006")
		}
		endOp := ""
		if row.endOp != "" {
			d, _ := time.Parse(time.RFC3339, row.endOp)
			endOp = d.Format("02.01.2006")
		}
		sectionStac := ""
		if row.sectionStac.Int64 > 0 {
			sectionStac = strconv.FormatInt(row.sectionStac.Int64, 10)
		}
		category := ""
		if row.category.Int64 > 0 {
			category = strconv.FormatInt(row.category.Int64, 10)
		}
		sectionReg := ""
		if row.sectionReg.Int64 > 0 {
			sectionReg = strconv.FormatInt(row.sectionReg.Int64, 10)
		}
		group := ""
		if row.group.Int64 > 0 {
			group = strconv.FormatInt(row.group.Int64, 10)
		}
		dateView := ""
		if row.dateView != "" {
			d, _ := time.Parse(time.RFC3339, row.dateView)
			dateView = d.Format("02.01.2006")
		}

		patientName, _ := utils.ToUTF8(row.patientName)
		address, _ := utils.ToUTF8(row.address)
		//mech, _ := utils.ToUTF8(row.mech)
		invalid, _ := utils.ToUTF8(row.invalid)
		//sick, _ := utils.ToUTF8(row.sick)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(cols, nRow), page.CellStyleBody(9))
		page.SetCellInt(1, nRow, row.patientId)
		page.SetCellStr(2, nRow, patientName)
		page.SetCellStr(3, nRow, bday)
		page.SetCellStr(4, nRow, startOp)
		page.SetCellStr(5, nRow, endOp)
		page.SetCellStr(6, nRow, row.mech)
		page.SetCellStr(7, nRow, sectionStac)
		page.SetCellStr(8, nRow, category)
		page.SetCellStr(9, nRow, sectionReg)
		page.SetCellStr(10, nRow, group)
		page.SetCellStr(11, nRow, dateView)
		page.SetCellStr(12, nRow, address)
		page.SetCellStr(13, nRow, invalid)
		page.SetCellStr(14, nRow, row.sick)

		nRow += 1
		sum += 1
	}
	sty = page.CellStyleTitle("left", "center", false, 9)
	page.Title(fmt.Sprintf("Всего: %v", sum), cellExcel(1, nRow), cellExcel(9, nRow), sty)
	nRow += 1
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(9, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(8, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func (r *ReportsJob) TakenForADNAccordingToClinical(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeDate) < 2 {
		return nil, errors.New(consts.ArrErrors[750])
	}

	sqlQuery := fmt.Sprintf(`select pid, fio, bd, uch, rd, diag from sps_pl_klin('%s','%s') order by fio`,
		p.Filters.RangeDate[0], p.Filters.RangeDate[1])

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		patientId   int
		patientName string
		bday        string
		section     int
		regDate     string
		diagnose    string
	}
	dateStart, _ := time.Parse("2006-01-02", p.Filters.RangeDate[0])
	dateEnd, _ := time.Parse("2006-01-02", p.Filters.RangeDate[1])

	cellExcel := excel.CellExcel
	toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	f.SetPageLayout(sheet,
		excelize.PageLayoutOrientation(excelize.OrientationPortrait),
		excelize.PageLayoutPaperSize(9),
	)
	page := excel.Page{f, sheet}

	sty := page.CellStyleTitle("center", "center", false, 9)
	collCount := 7
	// row 1
	page.Title(
		fmt.Sprintf("Список пациентов взятых на АДН по клиническим показаниям"),
		cellExcel(1, 1), cellExcel(collCount, 1), sty)
	page.Title(
		fmt.Sprintf("за период с %s по %s", dateStart.Format("02.01.2006"), dateEnd.Format("02.01.2006")),
		cellExcel(1, 2), cellExcel(collCount, 2), sty)
	// row 3
	f.SetRowHeight(sheet, 3, 27)
	f.SetCellStyle(sheet, cellExcel(1, 3), cellExcel(collCount, 3), page.CellStyleHeader(9))
	f.SetCellStr(sheet, cellExcel(1, 3), "")
	f.SetCellStr(sheet, cellExcel(2, 3), "Шифр")
	f.SetCellStr(sheet, cellExcel(3, 3), "Ф.И.О.")
	f.SetCellStr(sheet, cellExcel(4, 3), "Д.р.")
	f.SetCellStr(sheet, cellExcel(5, 3), "Участок")
	f.SetCellStr(sheet, cellExcel(6, 3), "Дата взятия")
	f.SetCellStr(sheet, cellExcel(7, 3), "Диагноз")

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 4.29)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(2), 7.29)
	f.SetColWidth(sheet, toCharStrConst(3), toCharStrConst(3), 28)
	f.SetColWidth(sheet, toCharStrConst(4), toCharStrConst(4), 10.14)
	f.SetColWidth(sheet, toCharStrConst(5), toCharStrConst(5), 10.14)
	f.SetColWidth(sheet, toCharStrConst(6), toCharStrConst(6), 10.29)
	f.SetColWidth(sheet, toCharStrConst(7), toCharStrConst(7), 10.14)

	nRow := 4
	nz := 1
	for rows.Next() {
		row := data{}
		rows.Scan(&row.patientId, &row.patientName, &row.bday, &row.section, &row.regDate, &row.diagnose)

		bd, _ := time.Parse(time.RFC3339, row.bday)
		bday := bd.Format("02.01.2006")

		bd, _ = time.Parse(time.RFC3339, row.regDate)
		regDate := bd.Format("02.01.2006")

		patientName, _ := utils.ToUTF8(row.patientName)

		f.SetCellStyle(sheet, cellExcel(1, nRow), cellExcel(collCount, nRow), page.CellStyleBody(9))
		page.SetCellInt(1, nRow, nz)
		page.SetCellInt(2, nRow, row.patientId)
		page.SetCellStr(3, nRow, patientName)
		page.SetCellStr(4, nRow, bday)
		page.SetCellInt(5, nRow, row.section)
		page.SetCellStr(6, nRow, regDate)
		page.SetCellStr(7, nRow, row.diagnose)

		nRow += 1
		nz += 1
	}
	sty = page.CellStyleTitle("left", "center", false, 9)
	page.Title("По данным отдела АСУ", cellExcel(1, nRow), cellExcel(collCount, nRow), sty)
	page.Title(fmt.Sprintf("Получено %s ", time.Now().Format("02.01.2006 15:04:05")), cellExcel(1, nRow+1), cellExcel(collCount, nRow+1), sty)

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func uklReportTypeOld(f *excelize.File, data types.UKLData, tx *sql.Tx) {
	sheet := f.GetSheetName(0)
	cellExcel := excel.CellExcel
	title := "Бланк  протокола по оценке результативности и качества труда медицинского персонала"
	title2 := "Критерии оценки качества медицинской помощи оказанной в амбулаторной помощи  врачом-психиатром участковым, врачом-психиатром детским участковым, врачом-психиатром подростковым участковым, врачом-психиатром (консультантом), врачом-психиатром детским (консультантом), врачом-психиатром подростковым (консультантом), врачом-сексологом, врачом-психотерапевтом, врачом-психиатром психиатрического отделения амбулаторного  принудительного лечения"

	page := excel.Page{f, sheet}

	collCount := 9

	// row 1
	sty := page.CellStyleTitle("center", "center", false, 11)
	page.Title(fmt.Sprintf(title),
		cellExcel(1, 1), cellExcel(collCount, 1), sty)

	sty = page.CellStyleHeader2(9)
	f.SetRowHeight(sheet, 2, 63)
	page.Title(fmt.Sprintf(title2),
		cellExcel(1, 2), cellExcel(collCount, 2), sty)

	sty = page.CellStyleTitle("left", "center", false, 9)
	f.SetRowHeight(sheet, 3, HEIGHT_ROW)
	page.Title(fmt.Sprintf("Вкладной лист к амбулаторной карте больного"),
		cellExcel(1, 3), cellExcel(collCount, 3), sty)

	patient, _ := models.Model.Patient.Get(int64(data.PatientId), tx)
	sty = page.CellStyleTitle("left", "center", false, 9)
	f.SetRowHeight(sheet, 4, HEIGHT_ROW)
	page.Title(fmt.Sprintf("Ф.И.О. пациента %s %s %s", patient.Lname, patient.Fname, patient.Sname),
		cellExcel(1, 4), cellExcel(collCount, 4), sty)

	doctor, _ := models.Model.User.Get(data.Doctor, tx)
	f.SetRowHeight(sheet, 5, HEIGHT_ROW)
	page.Title(fmt.Sprintf("Ф.И.О. врача  %s %s %s.", doctor.Lname, doctor.Fname, doctor.Sname),
		cellExcel(1, 5), cellExcel(collCount, 5), sty)

	f.SetRowHeight(sheet, 6, HEIGHT_ROW)
	page.Title(fmt.Sprintf("шифр %v", data.PatientId),
		cellExcel(1, 6), cellExcel(collCount, 6), sty)

	d, _ := time.Parse(time.RFC3339, data.Date1)
	f.SetRowHeight(sheet, 7, HEIGHT_ROW)
	page.Title(fmt.Sprintf("Дата заполнения %s", d.Format("02.01.2006")),
		cellExcel(1, 7), cellExcel(collCount, 7), sty)

	s1 := 0
	s2 := 0
	s3 := 0
	//table header
	f.SetCellStyle(sheet, cellExcel(1, 8), cellExcel(collCount, 9), page.CellStyleHeader(9))
	page.Range("Наименование критерия", cellExcel(1, 8), cellExcel(6, 9), 0)
	page.Range("Баллы", cellExcel(7, 8), cellExcel(9, 8), 0)
	f.SetCellStr(sheet, cellExcel(7, 9), "1 уровень")
	f.SetCellStr(sheet, cellExcel(8, 9), "2 уровень")
	f.SetCellStr(sheet, cellExcel(9, 9), "3 уровень")
	// end table header

	styTableTitle := page.CellStyleHeader(8)
	f.SetRowHeight(sheet, 10, HEIGHT_ROW)
	page.Range("1. Ведение медицинской документации - медицинской карты пациента, получающего медицинскую помощь в амбулаторных условиях:",
		cellExcel(1, 10), cellExcel(collCount, 10), styTableTitle)

	rowId := 11
	f.SetCellStyle(
		sheet,
		cellExcel(1, rowId),
		cellExcel(collCount, rowId+1),
		page.CellStyleBody(8),
	)
	page.Range("1.1 Заполнение всех разделов, предусмотренных амбулаторной картой",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_1)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_1)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_1)
	s1 += data.P1_1
	s2 += data.P2_1
	s3 += data.P3_1

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_2ROW)
	page.Range("1.2 Наличие информированного добровольного согласия на медицинское вмешательство",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_2)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_2)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_2)
	s1 += data.P1_2
	s2 += data.P2_2
	s3 += data.P3_2

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_ROW)
	page.Range("2. Первичный осмотр пациента и сроки оказания медицинской помощи:",
		cellExcel(1, rowId), cellExcel(collCount, rowId), styTableTitle)

	rowId += 1
	f.SetCellStyle(
		sheet,
		cellExcel(1, rowId),
		cellExcel(collCount, rowId+4),
		page.CellStyleBody(8),
	)
	f.SetRowHeight(sheet, rowId, HEIGHT_2ROW)
	page.Range("2.1 Оформление результатов первичного осмотра, включая данные анамнеза заболевания, записью в амбулаторной карте",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_3)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_3)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_3)
	s1 += data.P1_3
	s2 += data.P2_3
	s3 += data.P3_3

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_2ROW)
	page.Range("2.2 Установление предварительного диагноза лечащим врачом в ходе первичного приема пациента",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_4)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_4)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_4)
	s1 += data.P1_4
	s2 += data.P2_4
	s3 += data.P3_4

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_2ROW)
	page.Range("2.3 Формирование плана обследования пациента при первичном осмотре с учетом предварительного диагноза",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_5)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_5)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_5)
	s1 += data.P1_5
	s2 += data.P2_5
	s3 += data.P3_5

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_3ROW)
	page.Range("2.4 Формирование плана лечения при первичном осмотре с учетом предварительного диагноза, клинических проявлений заболевания, тяжести заболевания или состояния пациента",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_6)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_6)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_6)
	s1 += data.P1_6
	s2 += data.P2_6
	s3 += data.P3_6

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_4ROW)
	page.Range("2.5 Назначение лекарственных препаратов для медицинского применения с учетом инструкций по применению лекарственных препаратов, возраста пациента, пола пациента, тяжести заболевания, наличия осложнений основного заболевания (состояния) и сопутствующих заболеваний",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_7)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_7)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_7)
	s1 += data.P1_7
	s2 += data.P2_7
	s3 += data.P3_7

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_3ROW)
	page.Range("3. Установление клинического диагноза на основании данных анамнеза, осмотра, данных лабораторных, инструментальных и иных методов исследования, предусмотренных стандартами медицинской помощи, а также клинических рекомендаций (протоколов лечения) по вопросам оказания медицинской помощи:",
		cellExcel(1, rowId), cellExcel(collCount, rowId), styTableTitle)

	rowId += 1
	f.SetCellStyle(
		sheet,
		cellExcel(1, rowId),
		cellExcel(collCount, rowId+4),
		page.CellStyleBody(8),
	)
	f.SetRowHeight(sheet, rowId, HEIGHT_2ROW)
	page.Range("3.1 Оформление обоснования клинического диагноза соответствующей записью в амбулаторной карте",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_8)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_8)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_8)
	s1 += data.P1_8
	s2 += data.P2_8
	s3 += data.P3_8

	rowId += 1
	//f.SetRowHeight(sheet, rowId, HEIGHT_2ROW)
	page.Range("3.2 Установление клинического диагноза в течение 10 дней с момента обращения",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_9)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_9)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_9)
	s1 += data.P1_9
	s2 += data.P2_9
	s3 += data.P3_9

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_3ROW)
	page.Range("3.3 Проведение при затруднении установления клинического диагноза консилиума врачей с внесением соответствующей записи в амбулаторную карту с подписью зав. амбулаторно-поликлиническим отделением",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_10)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_10)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_10)
	s1 += data.P1_10
	s2 += data.P2_10
	s3 += data.P3_10

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_5ROW)
	page.Range("3.4 Внесение соответствующей записи в амбулаторную карту при наличии заболевания (состояния), требующего оказания мед. помощи в стационарных условиях, с указанием перечня рекомендуемых лабораторных и инструментальных методов исследований, а также оформление направления с указанием клинического диагноза при необходимости оказания мед.помощи в стационарных условиях в плановой форме",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_11)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_11)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_11)
	s1 += data.P1_11
	s2 += data.P2_11
	s3 += data.P3_11

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_4ROW)
	page.Range("3.5 Проведение коррекции плана обследования и плана лечения с учетом клинического диагноза, состояния пациента, особенностей течения заболевания, наличия сопутствующих заболеваний, осложнений заболевания и результатов проводимого лечения на основе стандартов мед. помощи и клинических рекомендаций",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_12)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_12)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_12)
	s1 += data.P1_12
	s2 += data.P2_12
	s3 += data.P3_12

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_ROW)
	page.Range("4. Назначение и выписывание лекарственных препаратов в соответствии с установленным порядком:",
		cellExcel(1, rowId), cellExcel(collCount, rowId), styTableTitle)

	rowId += 1
	f.SetCellStyle(
		sheet,
		cellExcel(1, rowId),
		cellExcel(collCount, rowId+4),
		page.CellStyleBody(8),
	)
	//f.SetRowHeight(sheet, rowId, HEIGHT_2ROW)
	page.Range("4.1 Оформление протокола решения врачебной комиссии",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_13)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_13)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_13)
	s1 += data.P1_13
	s2 += data.P2_13
	s3 += data.P3_13

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_3ROW)
	page.Range("4.2 Внесение записи в амбулаторную карту при назначении лекарственных препаратов для медицинского применения и применении медицинских изделий по решению врачебной комиссии",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_14)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_14)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_14)
	s1 += data.P1_14
	s2 += data.P2_14
	s3 += data.P3_14

	rowId += 1
	//f.SetRowHeight(sheet, rowId, HEIGHT_2ROW)
	page.Range("5. Проведение экспертизы временной нетрудоспособности в установленном порядке",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_15)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_15)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_15)
	s1 += data.P1_15
	s2 += data.P2_15
	s3 += data.P3_15

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_2ROW)
	page.Range("6. Осуществление диспансерного наблюдения в установленном порядке с соблюдением периодичности обследования и длительности диспансерного наблюдения",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), data.P1_16)
	f.SetCellInt(sheet, cellExcel(8, rowId), data.P2_16)
	f.SetCellInt(sheet, cellExcel(9, rowId), data.P3_16)
	s1 += data.P1_16
	s2 += data.P2_16
	s3 += data.P3_16

	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_ROW)
	page.Range("Итого:",
		cellExcel(1, rowId), cellExcel(6, rowId), 0)
	f.SetCellInt(sheet, cellExcel(7, rowId), s1)
	f.SetCellInt(sheet, cellExcel(8, rowId), s2)
	f.SetCellInt(sheet, cellExcel(9, rowId), s3)

	borderBottom, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	fontSize8, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: false,
			Size: 8,
		},
	})
	rowId += 2
	f.SetCellStyle(
		sheet,
		cellExcel(1, rowId),
		cellExcel(3, rowId+2),
		fontSize8,
	)
	f.SetRowHeight(sheet, rowId, HEIGHT_ROW)
	page.Range("Заведующий отделением", cellExcel(1, rowId), cellExcel(3, rowId), 0)
	page.Range("", cellExcel(4, rowId), cellExcel(5, rowId), borderBottom)
	page.Range("", cellExcel(7, rowId), cellExcel(8, rowId), borderBottom)
	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_ROW)
	page.Range("Заместитель главного врача", cellExcel(1, rowId), cellExcel(3, rowId), 0)
	page.Range("", cellExcel(4, rowId), cellExcel(5, rowId), borderBottom)
	page.Range("", cellExcel(7, rowId), cellExcel(8, rowId), borderBottom)
	rowId += 1
	f.SetRowHeight(sheet, rowId, HEIGHT_ROW)
	page.Range("Председатель ВК", cellExcel(1, rowId), cellExcel(3, rowId), 0)
	page.Range("", cellExcel(4, rowId), cellExcel(5, rowId), borderBottom)
	page.Range("", cellExcel(7, rowId), cellExcel(8, rowId), borderBottom)
}

func (r *ReportsJob) ProtocolUKL(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if p.Filters.Id == 0 {
		return nil, errors.New(consts.ArrErrors[753])
	}

	sqlQuery := fmt.Sprintf(`select nom_z,
p1_1, p1_2, p1_3, p1_4, p1_5, p1_6, p1_7, p1_8, p1_9, p1_10, p1_11, p1_12, p1_13, p1_14, p1_15, p1_16, p1_17, p1_18, p1_19, p1_20, p1_21, p1_22, p1_23, p1_24, p1_25, p1_26, p1_27, p1_28, p1_29, p1_30, p1_31, p1_32, p1_33, p1_34, p1_35,
p2_1, p2_2, p2_3, p2_4, p2_5, p2_6, p2_7, p2_8, p2_9, p2_10, p2_11, p2_12, p2_13, p2_14, p2_15, p2_16, p2_17, p2_18, p2_19, p2_20, p2_21, p2_22, p2_23, p2_24, p2_25, p2_26, p2_27, p2_28, p2_29, p2_30, p2_31, p2_32, p2_33, p2_34, p2_35,
p3_1, p3_2, p3_3, p3_4, p3_5, p3_6, p3_7, p3_8, p3_9, p3_10, p3_11, p3_12, p3_13, p3_14, p3_15, p3_16, p3_17, p3_18, p3_19, p3_20, p3_21, p3_22, p3_23, p3_24, p3_25, p3_26, p3_27, p3_28, p3_29, p3_30, p3_31, p3_32, p3_33, p3_34, p3_35,
NZ_REGISTRAT, p1_user, p2_user, p3_user, p1_date, p2_date, p3_date, dock, patient_id
from UKL where nom_z = ?`)

	row := tx.QueryRow(sqlQuery, p.Filters.Id)
	//dateStart, _ := time.Parse("2006-01-02", p.Filters.RangeDate[0])
	//dateEnd, _ := time.Parse("2006-01-02", p.Filters.RangeDate[1])

	data := types.UKLData{}
	err := row.Scan(
		&data.Id,
		&data.P1_1, &data.P1_2, &data.P1_3, &data.P1_4, &data.P1_5, &data.P1_6, &data.P1_7, &data.P1_8, &data.P1_9, &data.P1_10, &data.P1_11, &data.P1_12, &data.P1_13, &data.P1_14, &data.P1_15, &data.P1_16, &data.P1_17, &data.P1_18, &data.P1_19, &data.P1_20, &data.P1_21, &data.P1_22, &data.P1_23, &data.P1_24, &data.P1_25, &data.P1_26, &data.P1_27, &data.P1_28, &data.P1_29, &data.P1_30, &data.P1_31, &data.P1_32, &data.P1_33, &data.P1_34, &data.P1_35,
		&data.P2_1, &data.P2_2, &data.P2_3, &data.P2_4, &data.P2_5, &data.P2_6, &data.P2_7, &data.P2_8, &data.P2_9, &data.P2_10, &data.P2_11, &data.P2_12, &data.P2_13, &data.P2_14, &data.P2_15, &data.P2_16, &data.P2_17, &data.P2_18, &data.P2_19, &data.P2_20, &data.P2_21, &data.P2_22, &data.P2_23, &data.P2_24, &data.P2_25, &data.P2_26, &data.P2_27, &data.P2_28, &data.P2_29, &data.P2_30, &data.P2_31, &data.P2_32, &data.P2_33, &data.P2_34, &data.P2_35,
		&data.P3_1, &data.P3_2, &data.P3_3, &data.P3_4, &data.P3_5, &data.P3_6, &data.P3_7, &data.P3_8, &data.P3_9, &data.P3_10, &data.P3_11, &data.P3_12, &data.P3_13, &data.P3_14, &data.P3_15, &data.P3_16, &data.P3_17, &data.P3_18, &data.P3_19, &data.P3_20, &data.P3_21, &data.P3_22, &data.P3_23, &data.P3_24, &data.P3_25, &data.P3_26, &data.P3_27, &data.P3_28, &data.P3_29, &data.P3_30, &data.P3_31, &data.P3_32, &data.P3_33, &data.P3_34, &data.P3_35,
		&data.RegistratId, &data.User1, &data.User2, &data.User3, &data.Date1, &data.Date2, &data.Date3, &data.Doctor, &data.PatientId,
	)
	if err != nil {
		return nil, err
	}
	//unit := 0
	//if data.VisitId != 0 {
	//	unit, _ = getVisitById(data.VisitId, tx)
	//}

	//toCharStrConst := excel.ToCharStrConst

	f := excel.CreateFile()
	sheet := f.GetSheetName(0)
	f.SetPageLayout(sheet,
		excelize.PageLayoutOrientation(excelize.OrientationPortrait),
		excelize.PageLayoutPaperSize(9),
	)
	f.SetPageMargins(sheet,
		excelize.PageMarginBottom(0.5),
		excelize.PageMarginFooter(1.0),
		excelize.PageMarginHeader(1.0),
		excelize.PageMarginLeft(0.5),
		excelize.PageMarginRight(0),
		excelize.PageMarginTop(0.5),
	)

	uklReportTypeOld(f, data, tx)
	//if unit == 2 {
	//	uklReportTypeOld(f, data, tx)
	//} else if unit == 4 {
	//
	//} else {
	//	uklReportTypeOld(f, data, tx)
	//}

	buf, _ := f.WriteToBuffer()
	return buf, nil
}

func (r *ReportsJob) form39GenerateData(d1 string, d2 string) {
	_, tx := CreateTx()
	defer tx.Commit()
	//maxDay, err := models.Model.Spr.GetParam("max_day", tx)
	//if err != nil {
	//	ERROR.Println(err)
	//	return
	//}
	//dateEnd, err := time.Parse(consts.DATE_FORMAT_INPUT, d2)
	//if err != nil {
	//	ERROR.Println(err)
	//	return
	//}
	//now := time.Now()
	//if (now.Day() <= maxDay.ParamI && dateEnd.Month() == now.Month()-1 && dateEnd.Year() == now.Year()) ||
	//	(now.Month() == dateEnd.Month() && now.Year() == dateEnd.Year()) {
	err := form39GenerateData(d1, d2, tx)
	if err != nil {
		ERROR.Println(err)
	}
	//}
}

func (r *ReportsJob) Form39General(p reportParams, tx *sql.Tx) (*bytes.Buffer, error) {
	if len(p.Filters.RangeDate) < 2 {
		return nil, errors.New(consts.ArrErrors[750])
	}
	dateEnd, _ := time.Parse(consts.DATE_FORMAT_INPUT, p.Filters.RangeDate[1])
	dateStart := time.Date(dateEnd.Year(), dateEnd.Month(), 1, 0, 0, 0, 0, dateEnd.Location())
	p.Filters.RangeDate[0] = dateStart.Format(consts.DATE_FORMAT_INPUT)
	p.Filters.RangeDate[1] = dateEnd.Format(consts.DATE_FORMAT_INPUT)

	cacheName := fmt.Sprintf("Form39General_%s_%s_%v", p.Filters.RangeDate[0], p.Filters.RangeDate[1], p.Filters.Unit)
	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		return item.(*bytes.Buffer), nil
	}

	r.form39GenerateData(dateStart.Format(consts.DATE_FORMAT_INPUT), dateEnd.Format(consts.DATE_FORMAT_INPUT))

	sqlQuery := fmt.Sprintf(`select name_doct, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, p18, p19, p20 from F39_ARHIV_SVOD('%s', '%s', %v)`,
		p.Filters.RangeDate[0], p.Filters.RangeDate[1], p.Filters.Unit)
	INFO.Println(sqlQuery)
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	type data struct {
		doctorName string
		p1         int
		p2         int
		p3         int
		p4         int
		p5         int
		p6         int
		p7         int
		p8         int
		p9         int
		p10        int
		p11        int
		p12        int
		p13        int
		p14        int
		p15        int
		p16        int
		p17        int
		p18        int
		p19        int
		p20        int
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
	collCount := 21
	f.SetPageMargins(sheet,
		excelize.PageMarginBottom(0.5),
		excelize.PageMarginFooter(1.0),
		excelize.PageMarginHeader(1.0),
		excelize.PageMarginLeft(0.5),
		excelize.PageMarginRight(0),
		excelize.PageMarginTop(0.5),
	)

	//dateStart, _ := time.Parse("2006-01-02", p.Filters.RangeDate[0])
	//dateEnd, _ := time.Parse("2006-01-02", p.Filters.RangeDate[1])

	f.SetColWidth(sheet, toCharStrConst(1), toCharStrConst(1), 30)
	f.SetColWidth(sheet, toCharStrConst(2), toCharStrConst(collCount), 5.5)

	styCC := page.CellStyleTitle("center", "center", false, 9)
	styLC := page.CellStyleTitle("left", "center", false, 9)
	//styRC := page.CellStyleTitle("right", "center", false, 9)
	// row 1
	rowId := 1
	page.Title("Министерство здравоохранения", cellExcel(1, rowId), cellExcel(1, rowId), styCC)
	page.Title("СВОДНАЯ ВЕДОМОСТЬ", cellExcel(2, rowId), cellExcel(14, rowId+3), styCC)
	page.Title(fmt.Sprintf("Дата создания списка: %s", time.Now().Format(consts.DATE_FORMAT_RU)), cellExcel(15, rowId), cellExcel(collCount, rowId), styLC)
	rowId += 1
	page.Title("Российской Федерации", cellExcel(1, rowId), cellExcel(1, rowId), styCC)
	page.Title("Медицинская документация", cellExcel(16, rowId), cellExcel(collCount, rowId), styLC)
	rowId += 1
	page.Title("БУЗОО КПБ им. Солодникова Н.Н.", cellExcel(1, rowId), cellExcel(1, rowId), styCC)
	page.Title("форма № 039/у-02", cellExcel(16, rowId), cellExcel(collCount, rowId), styLC)
	rowId += 1
	page.Title("утверждена приказом Минздравом России", cellExcel(16, rowId), cellExcel(collCount, rowId), styLC)
	rowId += 1
	unitName := "Взрослый диспансер"
	switch p.Filters.Unit {
	case 1:
		unitName = "Взрослый диспансер"
	case 2:
		unitName = "Психотерапия"
	case 4:
		unitName = "Суицидология"
	case 8:
		unitName = "ОИЛС"
	case 16:
		unitName = "Детский диспансер"
	case 1024:
		unitName = "Специалисты"
	case 512:
		unitName = "Село"
	case 2048:
		unitName = "АПЛ"
	}
	page.Title(unitName, cellExcel(1, rowId), cellExcel(1, rowId), styLC)
	page.Title("учета врачебных посещений в амбулаторно-поликлинических учреждениях, на дому", cellExcel(2, rowId), cellExcel(14, rowId), styCC)
	page.Title("от 30.12.2002 №413", cellExcel(16, rowId), cellExcel(collCount, rowId), styLC)
	rowId += 1
	page.Title(fmt.Sprintf("с %s - %s", dateStart.Format(consts.DATE_FORMAT_RU), dateEnd.Format(consts.DATE_FORMAT_RU)), cellExcel(2, rowId), cellExcel(14, rowId), styCC)
	//table header
	styHeader := page.CellStyleHeader(9)
	rowId += 1
	page.Title("Ф.И.О. врач", cellExcel(1, rowId), cellExcel(1, rowId+2), styHeader)
	page.Title("ВСЕГО", cellExcel(2, rowId), cellExcel(2, rowId+2), styHeader)
	page.Title("Число посещений в поликлинике", cellExcel(3, rowId), cellExcel(4, rowId+1), styHeader)
	page.Title("В том числе в возрасте (из гарфы2)", cellExcel(5, rowId), cellExcel(6, rowId+1), styHeader)
	page.Title("Из общего числа посещений в поликлинике по поводу заболеваний", cellExcel(7, rowId), cellExcel(9, rowId+1), styHeader)
	page.Title("Из гр. 3 - профилактических", cellExcel(10, rowId), cellExcel(10, rowId+2), styHeader)
	page.Title("Число посещений на дому всего", cellExcel(11, rowId), cellExcel(11, rowId+2), styHeader)
	page.Title("из общего числа посещений на дому", cellExcel(12, rowId), cellExcel(17, rowId), styHeader)
	page.Title("Число посещений по видам оплаты", cellExcel(18, rowId), cellExcel(21, rowId), styHeader)
	rowId += 1
	f.SetRowHeight(sheet, rowId, 38)
	page.Title("по поводу заболеваний", cellExcel(12, rowId), cellExcel(15, rowId), styHeader)
	page.Title("из числа профилактических", cellExcel(16, rowId), cellExcel(17, rowId), styHeader)
	page.Title("ОМС", cellExcel(18, rowId), cellExcel(18, rowId+1), styHeader)
	page.Title("Бюджет", cellExcel(19, rowId), cellExcel(19, rowId+1), styHeader)
	page.Title("Платные", cellExcel(20, rowId), cellExcel(20, rowId+1), styHeader)
	page.Title("ДМС", cellExcel(21, rowId), cellExcel(21, rowId+1), styHeader)
	rowId += 1
	f.SetRowHeight(sheet, rowId, 38)
	page.Title("Всего", cellExcel(3, rowId), cellExcel(3, rowId), styHeader)
	page.Title("Сельских", cellExcel(4, rowId), cellExcel(4, rowId), styHeader)
	page.Title("0-17 лет", cellExcel(5, rowId), cellExcel(5, rowId), styHeader)
	page.Title("60 лет и старше", cellExcel(6, rowId), cellExcel(6, rowId), styHeader)
	page.Title("Всего", cellExcel(7, rowId), cellExcel(7, rowId), styHeader)
	page.Title("0-17 лет", cellExcel(8, rowId), cellExcel(8, rowId), styHeader)
	page.Title("60 лет старше", cellExcel(9, rowId), cellExcel(9, rowId), styHeader)
	page.Title("Всего", cellExcel(12, rowId), cellExcel(12, rowId), styHeader)
	page.Title("0-17 лет", cellExcel(13, rowId), cellExcel(13, rowId), styHeader)
	page.Title("из них 0-1 год", cellExcel(14, rowId), cellExcel(14, rowId), styHeader)
	page.Title("60 лет и старше", cellExcel(15, rowId), cellExcel(15, rowId), styHeader)
	page.Title("0-17 лет", cellExcel(16, rowId), cellExcel(16, rowId), styHeader)
	page.Title("0-1 год", cellExcel(17, rowId), cellExcel(17, rowId), styHeader)

	styBodyL, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: false,
			Vertical: "left",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
		},
		Font: &excelize.Font{
			Size: 9,
		},
	})
	styBodyR, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: false,
			Vertical: "right",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
		},
		Font: &excelize.Font{
			Size: 9,
		},
	})
	rowId += 1
	for rows.Next() {
		row := data{}
		err := rows.Scan(&row.doctorName, &row.p1, &row.p2, &row.p3, &row.p4, &row.p5, &row.p6, &row.p7, &row.p8, &row.p9,
			&row.p10, &row.p11, &row.p12, &row.p13, &row.p14, &row.p15, &row.p16, &row.p17, &row.p18, &row.p19, &row.p20)
		if err != nil {
			ERROR.Println(err)
			return nil, err
		}

		doctorName, _ := utils.ToUTF8(row.doctorName)

		f.SetCellStyle(sheet, cellExcel(1, rowId), cellExcel(1, rowId), styBodyL)
		f.SetCellStyle(sheet, cellExcel(2, rowId), cellExcel(21, rowId), styBodyR)
		page.SetCellStr(1, rowId, doctorName)
		page.SetCellInt(2, rowId, row.p1)
		page.SetCellInt(3, rowId, row.p2)
		page.SetCellInt(4, rowId, row.p3)
		page.SetCellInt(5, rowId, row.p4)
		page.SetCellInt(6, rowId, row.p5)
		page.SetCellInt(7, rowId, row.p6)
		page.SetCellInt(8, rowId, row.p7)
		page.SetCellInt(9, rowId, row.p8)
		page.SetCellInt(10, rowId, row.p9)
		page.SetCellInt(11, rowId, row.p10)
		page.SetCellInt(12, rowId, row.p11)
		page.SetCellInt(13, rowId, row.p12)
		page.SetCellInt(14, rowId, row.p13)
		page.SetCellInt(15, rowId, row.p14)
		page.SetCellInt(16, rowId, row.p15)
		page.SetCellInt(17, rowId, row.p16)
		page.SetCellInt(18, rowId, row.p17)
		page.SetCellInt(19, rowId, row.p18)
		page.SetCellInt(20, rowId, row.p19)
		page.SetCellInt(21, rowId, 0)

		rowId += 1
	}

	buf, _ := f.WriteToBuffer()
	cache.AppCache.Set(cacheName, buf, time.Minute)
	return buf, nil
}
