package controller

import (
	"fmt"
	"strings"
	"time"

	"pmain2/internal/consts"
	"pmain2/internal/database"
	"pmain2/internal/models"
	"pmain2/internal/types"
	"pmain2/pkg/cache"
	"pmain2/pkg/utils"
)

var (
	cachePat = cache.CreateCache(time.Minute, time.Minute)
)

type patient struct{}

func initPatientController() *patient {
	return &patient{}
}

func (p *patient) FindByFio(lname, fname, sname string) (*[]models.Patient, error) {
	cacheName := lname + " " + fname + " " + sname

	item, ok := cachePat.Get(cacheName)
	if ok {
		res := item.(*[]models.Patient)
		return res, nil
	}

	model := models.Model.Patient
	lname, _ = utils.ToWin1251(lname)
	fname, _ = utils.ToWin1251(fname)
	sname, _ = utils.ToWin1251(sname)
	data, err := model.FindByFIO(lname, fname, sname)
	if err != nil {
		return nil, err
	}

	cachePat.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) FindById(id int) (*models.Patient, error) {

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Patient
	data, err := model.Get(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *patient) FindUchet(id int, isCache bool) (*[]models.FindUchetS, error) {
	cacheName := fmt.Sprintf("find_uchet_%v", id)
	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		return item.(*[]models.FindUchetS), nil
	}
	model := models.Model.Patient
	data, err := model.FindUchet(id, 1000, 0)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) FindLastUchet(id int, isCache bool) (*models.FindUchetS, error) {
	cacheName := fmt.Sprintf("find_last_uchet_%v", id)
	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		return item.(*models.FindUchetS), nil
	}
	model := models.Model.Patient
	data, err := model.FindLastUchet(id)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 30)
	return data, nil
}

func (p *patient) HistoryVisits(id int, isCache bool) (*[]models.HistoryVisit, error) {
	cacheName := fmt.Sprintf("disp_history_Visit_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]models.HistoryVisit), nil
		}
	}
	model := models.Model.Patient
	data, err := model.HistoryVisits(id)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) HistoryHospital(id int) (*[]models.HistoryHospital, error) {
	cacheName := fmt.Sprintf("disp_history_hospital_%v", id)
	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		return item.(*[]models.HistoryHospital), nil
	}
	model := models.Model.Patient
	data, err := model.HistoryHospital(id)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewVisit(visit *types.NewVisit) (int, error) {
	fmt.Println(*visit)
	visit.Normalize()
	model := models.Model.Patient
	lastUchet, err := model.FindLastUchet(visit.PatientId)
	if err != nil {
		return 100, err
	}

	//-проверить что пациент не мертв или это работа с документами
	if (lastUchet != nil && lastUchet.Reason == consts.EXIT_REAS_DEAD) && visit.Visit&consts.VISIT_WORK_WITH_DOCUMENTS == 0 {
		return 101, nil
	}

	//-в этот день не было посещений
	isVisisted, err := model.IsVisited(visit)
	if err != nil {
		return 102, err
	}

	if isVisisted {
		return 202, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return 20, err
	}
	tx, err := conn.DB.Begin()
	if err != nil {
		return 21, err
	}

	//Обрезаем до 10, т.к. в посещениях длина диагноза 10
	if len(visit.Diagnose) > 10 {
		visit.Diagnose = visit.Diagnose[0:10]
	}

	model = models.Model.Patient
	_, err = model.NewVisit(*visit, tx)
	if err != nil {
		tx.Rollback()
		return 200, err
	}
	if visit.SRC >= 0 {
		_, err = model.NewSRC(&types.NewSRC{
			PatientId: visit.PatientId,
			DateAdd:   visit.Date,
			DockId:    visit.DockId,
			Unit:      visit.Unit,
			Zakl:      visit.SRC,
		}, tx)
		if err != nil {
			tx.Rollback()
			return 201, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) NewProf(visit *types.NewProf) (int, error) {
	fmt.Println(*visit)
	visit.Normalize()
	model := models.Model.Patient

	conn, err := database.Connect()
	if err != nil {
		return 20, err
	}
	tx, err := conn.DB.Begin()
	if err != nil {
		return 21, err
	}

	if visit.Count == 0 {
		return 203, nil
	}

	model = models.Model.Patient
	for i := 0; i < visit.Count; i++ {
		_, err = model.NewProf(*visit, tx)
		if err != nil {
			tx.Rollback()
			return 200, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) NewReg(reg *types.NewRegister) (int, error) {
	pat, err := p.FindById(reg.PatientId)
	if err != nil {
		return -1, err
	}

	sprModel := models.Model.Spr
	patientModel := models.Model.Patient

	if reg.Reason == "001" {
		if len(pat.Address) < 10 {
			return 301, nil
		}
	}

	lastReg, err := p.FindLastUchet(reg.PatientId, false)
	if err != nil {
		return -1, err
	}

	isClose, err := sprModel.IsClosedSection(lastReg.Section)
	if err != nil {
		return -1, err
	}
	if isClose {
		return 303, nil
	}

	if reg.Diagnose == "" {
		reg.Diagnose = lastReg.Diagnose
	}

	if reg.Section == 0 {
		reg.Section = lastReg.Section
	}

	if reg.Category == 0 {
		reg.Category = lastReg.Category
	}

	countJudgment, err := patientModel.GetCountJudgment(reg.PatientId)
	if err != nil {
		return -1, err
	}
	if reg.Category == 7 && countJudgment == 0 {
		return 304, nil
	}

	if lastReg.Reason == consts.EXIT_REAS_DEAD {
		return 305, nil
	}

	if reg.Section < 10 && reg.Reason == consts.REAS_NEW {
		return 306, nil
	}

	if (reg.Category == 7 || reg.Category == 8) && (reg.Section < 18 || reg.Section > 19) && reg.Section < 130 {
		return 307, nil
	}

	if reg.Reason == consts.REAS_SWITCH_CATEG_TO_AMBULANC && reg.Category == 10 {
		return 308, nil
	}

	if reg.Reason == consts.REAS_SWITCH_CATEG_TO_CONSULTANT && reg.Category != 10 {
		return 309, nil
	}

	if reg.Reason == consts.REAS_SWITCH_CATEG_GROUP && (reg.Category == 10 || reg.Category == lastReg.Category) {
		return 310, nil
	}

	if (reg.Reason == consts.REAS_SWITCH_CATEG_TO_AMBULANC && lastReg.Category > 0 && lastReg.Category < 9) ||
		(reg.Reason == consts.REAS_SWITCH_CATEG_TO_CONSULTANT && lastReg.Category == 10) {
		return 311, nil
	}

	if reg.Reason == consts.REAS_EXIT {
		if reg.ExitReason == "" {
			return 316, nil
		}
		reg.Reason = reg.ExitReason
	}

	if reg.Reason == consts.EXIT_REAS_NO_PSIH_DIAG {
		//TODO для всех prava <> 2147483647
		return 312, nil
	}

	isClose, err = sprModel.IsClosedSection(reg.Section)
	if err != nil {
		return -1, err
	}
	if isClose {
		return 313, nil
	}

	if reg.Reason[0] == 'S' {
		inHospital, err := patientModel.IsInHospital(reg.PatientId)
		if err != nil {
			return -1, err
		}
		if inHospital {
			return 314, nil
		}
	}

	regDate, err := time.Parse("2006-01-02", reg.Date)
	if err != nil {
		return -1, err
	}
	countReg, err := patientModel.GetCountRegDataInDate(reg.PatientId, reg.Section, regDate)
	if err != nil {
		return -1, err
	}
	if countReg > 0 {
		return 315, nil
	}

	lastRegDate, err := time.Parse("2006-01-02", lastReg.Date)
	if err != nil {
		return -1, err
	}

	if lastRegDate.Sub(regDate) > 0 {
		return 360, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return 20, err
	}
	tx, err := conn.DB.Begin()
	if err != nil {
		return 21, err
	}

	defer tx.Rollback()
	_, err = patientModel.InsertReg(*reg, tx)
	if err != nil {
		return 350, err
	}
	if reg.Reason == consts.REAS_NEW {
		_, err = patientModel.UpdPatientVisible(reg.PatientId, 0, tx)
		if err != nil {
			return 351, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 350, err
	}

	return 0, nil
}

func (p *patient) NewRegisterTransfer(reg *types.NewRegisterTransfer) (int, error) {
	//pat, err := p.FindById(reg.PatientId)
	//if err != nil {
	//	return -1, err
	//}
	sprModel := models.Model.Spr
	patientModel := models.Model.Patient

	lastReg, err := p.FindLastUchet(reg.PatientId, false)
	if err != nil {
		return -1, err
	}

	isClose, err := sprModel.IsClosedSection(lastReg.Section)
	if err != nil {
		return -1, err
	}
	if isClose {
		return 303, nil
	}

	isClose, err = sprModel.IsClosedSection(reg.Section)
	if err != nil {
		return -1, err
	}
	if isClose {
		return 303, nil
	}

	if ((lastReg.Section == 19 && reg.Section == 18) || (lastReg.Section == 18 && reg.Section == 19)) && reg.Category < 7 {
		return 370, nil
	}
	if ((reg.Section == 18) || (reg.Section == 19)) && reg.Category < 7 {
		return 370, nil
	}
	if ((lastReg.Section == 18) || (lastReg.Section == 19)) && reg.Category > 6 {
		return 370, nil
	}

	reasonPrev := consts.REAS_FROM
	reasonNext := consts.REAS_TO

	if (lastReg.Section == 481 && reg.Section == 480) ||
		(lastReg.Section == 591 && reg.Section == 590) ||
		(lastReg.Section == 661 && reg.Section == 660) {
		reasonPrev = consts.REAS_FROM
		reasonNext = consts.REAS_TO
	} else {
		if lastReg.Section >= 400 && reg.Section < 400 {
			reasonPrev = consts.EXIT_REAS_EXIT
			reasonNext = consts.REAS_NEW
		}
		if lastReg.Section < 400 && reg.Section >= 400 {
			reasonPrev = consts.EXIT_REAS_EXIT
			reasonNext = consts.REAS_NEW
		}
	}
	regDate, err := time.Parse("2006-01-02", reg.Date)
	if err != nil {
		return -1, err
	}

	lastRegDate, err := time.Parse("2006-01-02", lastReg.Date)
	if err != nil {
		return -1, err
	}

	if lastRegDate.Sub(regDate) > 0 {
		return 360, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return 20, err
	}
	tx, err := conn.DB.Begin()
	if err != nil {
		return 21, err
	}

	defer tx.Rollback()
	if err != nil {
		return -1, err
	}
	regPrev := types.NewRegister{
		PatientId:  reg.PatientId,
		Reason:     reasonPrev,
		ExitReason: "",
		Section:    lastReg.Section,
		Category:   lastReg.Category,
		Diagnose:   lastReg.Diagnose,
		Date:       reg.Date,
		DockId:     reg.DockId,
	}
	_, err = patientModel.InsertReg(regPrev, tx)
	if err != nil {
		return 350, err
	}
	regNext := types.NewRegister{
		PatientId:  reg.PatientId,
		Reason:     reasonNext,
		ExitReason: "",
		Section:    reg.Section,
		Category:   reg.Category,
		Diagnose:   lastReg.Diagnose,
		Date:       reg.Date,
		DockId:     reg.DockId,
	}
	_, err = patientModel.InsertReg(regNext, tx)
	if err != nil {
		return 350, err
	}
	err = tx.Commit()
	if err != nil {
		return 350, err
	}

	return 0, nil
}

func (p *patient) HistorySindrom(id int, isCache bool) (*[]models.HistorySindrom, error) {
	cacheName := fmt.Sprintf("disp_sindrom_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]models.HistorySindrom), nil
		}
	}
	model := models.Model.Patient
	data, err := model.HistorySindrom(id)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewSindrom(sindrom *types.Sindrom) (int, error) {

	conn, err := database.Connect()
	if err != nil {
		return 20, err
	}
	tx, err := conn.DB.Begin()
	if err != nil {
		return 21, err
	}

	model := models.Model.Patient

	history, err := model.HistorySindrom(sindrom.PatientId)
	if err != nil {
		return -1, err
	}

	isSindrom := strings.Contains(sindrom.Diagnose, "F")

	count := 0
	for _, row := range *history {
		isSindromRow := strings.Contains(row.Diagnose, "F")
		if (isSindrom && isSindromRow) || (!isSindrom && !isSindromRow) {
			count += 1
		}
	}

	if isSindrom && count >= 4 {
		return 380, nil
	}
	if !isSindrom && count >= 4 {
		return 381, nil
	}

	_, err = model.NewSindrom(*sindrom, tx)
	if err != nil {
		tx.Rollback()
		return 200, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) RemoveSindrom(sindrom *types.Sindrom) (int, error) {

	conn, err := database.Connect()
	if err != nil {
		return 20, err
	}
	tx, err := conn.DB.Begin()
	if err != nil {
		return 21, err
	}

	model := models.Model.Patient

	_, err = model.RemoveSindrom(*sindrom, tx)
	if err != nil {
		tx.Rollback()
		return 200, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) FindInvalid(id int, isCache bool) (*[]models.FindInvalid, error) {
	cacheName := fmt.Sprintf("find_invalid_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]models.FindInvalid), nil
		}
	}
	model := models.Model.Patient
	data, err := model.FindInvalid(id)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}
