package types

type Find struct {
	Name string `json:"name"`
}

type FindI struct {
	Name int64 `json:"name"`
}

type Spr struct {
	Code  string `json:"code"`
	Param string `json:"param"`
	Name  string `json:"name"`
}

type Patient struct {
	Id             int64  `json:"id"`
	Lname          string `json:"lname"`
	Fname          string `json:"fname"`
	Sname          string `json:"sname"`
	Bday           string `json:"bday"`
	Visibility     int    `json:"visibility"`
	Sex            string `json:"sex"`
	Snils          string `json:"snils"`
	Address        string `json:"address"`
	PassportSeries string `json:"passportSeries"`
	PassportNumber int    `json:"passportNumber"`
	Works          int    `json:"works"`
	Republic       int    `json:"republic"`
	Region         int    `json:"region"`
	District       int    `json:"district"`
	Area           int    `json:"area"`
	Street         int    `json:"street"`
	House          string `json:"house"`
	Build          string `json:"build"`
	Flat           string `json:"flat"`
	Domicile       int    `json:"domicile"`
}

type NewPatient struct {
	PatientId int64  `json:"patientId"`
	Lname     string `json:"lname"`
	Fname     string `json:"fname"`
	Sname     string `json:"sname"`
	Bday      string `json:"bday"`
	IsAnonim  bool   `json:"isAnonim"`
	Sex       string `json:"sex"`
	UserId    int    `json:"userId"`
	IsForced  bool   `json:"isForced"`
}

type Visit struct {
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

type NewVisit struct {
	Visit       int    `json:"visit"`
	Uch         int    `json:"uch"`
	Unit        int    `json:"unit"`
	Home        bool   `json:"home"`
	Diagnose    string `json:"diagnose"`
	Date        string `json:"date"`
	PatientId   int64  `json:"patientId"`
	PatientBDay string `json:"patientBDay"`
	DockId      int    `json:"dockId"`
	SRC         int    `json:"src"`
}

type NewProf struct {
	Count  int    `json:"count"`
	Date   string `json:"date"`
	Unit   int    `json:"unit"`
	DockId int    `json:"dockId"`
	Uch    int    `json:"uch"`
	Home   bool   `json:"home"`
}

type Sindrom struct {
	Id        int    `json:"id"`
	PatientId int    `json:"patientId"`
	Diagnose  string `json:"diagnose"`
	DoctId    int    `json:"doctId"`
}

type NewSRC struct {
	PatientId int64  `json:"patientId"`
	DateAdd   string `json:"dateAdd"`
	DockId    int    `json:"dockId"`
	Unit      int    `json:"unit"`
	Zakl      int    `json:"zakl"`
}

type NewRegister struct {
	PatientId  int64  `json:"patientId"`
	Reason     string `json:"reason"`
	ExitReason string `json:"exitReason"`
	Section    int    `json:"section"`
	Category   int    `json:"category"`
	Diagnose   string `json:"diagnose"`
	Date       string `json:"date"`
	DockId     int    `json:"dockId"`
}
type NewRegisterTransfer struct {
	PatientId int64  `json:"patientId"`
	Category  int    `json:"category"`
	Section   int    `json:"section"`
	Date      string `json:"date"`
	DockId    int    `json:"dockId"`
}

type HttpResponse struct {
	Success bool        `json:"success"`
	Error   int         `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type NewInvalid struct {
	DoctId       int    `json:"doctId"`
	PatientId    int64  `json:"patientId"`
	DateStart    string `json:"date_start"`
	DateEnd      string `json:"date_end"`
	DateDocument string `json:"date_document"`
	Reason       string `json:"reason"`
	Kind         string `json:"kind"`
	Anomaly      string `json:"anomaly"`
	Limit        string `json:"limit"`
	IsInfinity   bool   `json:"isInfinity"`
}

type FindCustody struct {
	DateStart string `json:"dateStart"`
	DateEnd   string `json:"dateEnd"`
	Who       string `json:"who"`
}

type NewCustody struct {
	PatientId int64  `json:"patientId"`
	DoctId    int    `json:"doctId"`
	Custody   string `json:"custody"`
	DateStart string `json:"dateStart"`
	DateEnd   string `json:"dateEnd"`
}

type FindVaccination struct {
	Date        string `json:"date"`
	Vaccination string `json:"vaccination"`
	Number      string `json:"number"`
	Series      string `json:"series"`
	Result      string `json:"result"`
	Detached    string `json:"detached"`
}

type FindInfection struct {
	Date     string `json:"date"`
	Diagnose string `json:"diagnose"`
}

type ST22 struct {
	Id        int    `json:"id"`
	PatientId int64  `json:"patientId"`
	DateStart string `json:"dateStart"`
	DateEnd   string `json:"dateEnd"`
	Section   int    `json:"section"`
	Part      int    `json:"part"`
	InsWho    int    `json:"insWho"`
	InsDate   string `json:"insDate"`
}

type SprUchN struct {
	Id      int     `json:"id"`
	Section int     `json:"section"`
	Name    string  `json:"name"`
	Plan    float64 `json:"plan"`
	Hour    float64 `json:"hour"`
	Spec    string  `json:"spec"`
	Unit    int     `json:"unit"`
}

type SOD struct {
	PatientId int64  `json:"patientId"`
	Date      string `json:"date"`
	Section   int    `json:"section"`
	Part      int    `json:"part"`
}

type OOD struct {
	PatientId    int64 `json:"patientId"`
	Danger       int   `json:"danger"`
	Syndrome     int   `json:"syndrome"`
	Difficulties int   `json:"difficulties"`
	Attitude     int   `json:"attitude"`
	UserId       int64 `json:"userId"`
}

type FindSection29 struct {
	DateStart string `json:"dateStart"`
	Diagnose  string `json:"diagnose"`
	DateEnd   string `json:"dateEnd"`
	Section   string `json:"section"`
}

type LocationDoctor struct {
	Section int    `json:"section"`
	Spec    string `json:"spec"`
	DoctId  int    `json:"doctId"`
	Lname   string `json:"lname"`
	Fname   string `json:"fname"`
	Sname   string `json:"sname"`
	Unit    int    `json:"unit"`
}

//administration

type DoctorBySection struct {
	Section  int `json:"section"`
	DoctorId int `json:"doctorId"`
}

type NewDoctorLocation struct {
	Date string            `json:"date"`
	Unit int               `json:"unit"`
	Data []DoctorBySection `json:"data"`
}

type NewDoctorLeadSection struct {
	Month int               `json:"month"`
	Year  int               `json:"year"`
	Unit  int               `json:"unit"`
	Data  []DoctorBySection `json:"data"`
}

type FindDoctorLead struct {
	Unit  int `json:"unit"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type ForcedM struct {
	Id        int    `json:"id"`
	DateStart string `json:"dateStart"`
	DateEnd   string `json:"dateEnd"`
	Watch     string `json:"watch"`
	Mechanism string `json:"mechanism"`
	State     string `json:"state"`
	Number    int    `json:"number"`
}

type Forced struct {
	Id                  int    `json:"id"`
	PatientId           int    `json:"patientId"`
	Number              int    `json:"number"`
	CourtDate           string `json:"courtDate"`
	CourtConclusionDate string `json:"courtConclusionDate"`
	CourtId             int    `json:"courtId"`
	TypeCrimeId         int    `json:"typeCrimeId"`
	ViewId              int    `json:"viewId"`
	ForcedP             int    `json:"forcedP"`
	Sick                int    `json:"sick"`
	Mechanism           int    `json:"mechanism"`
	DoctorId1           int    `json:"doctorId1"`
	DoctorId2           int    `json:"doctorId2"`
	DateView            string `json:"dateView"`
	ConclusionId        int    `json:"conclusionId"`
	ActNumber           int    `json:"actNumber"`
	ActDate             string `json:"actDate"`
	DateEnd             string `json:"dateEnd"`
	TypeId              int    `json:"typeId"`
	UserId              int    `json:"userId"`
}

type ViewedM struct {
	Id                  int    `json:"id"`
	ViewDate            string `json:"viewDate"`
	DoctorName1         string `json:"doctorName1"`
	DoctorName2         string `json:"doctorName2"`
	Conclusion          string `json:"conclusion"`
	ActNumber           int    `json:"actNumber"`
	ActDate             string `json:"actDate"`
	View                string `json:"view"`
	CourtDate           string `json:"courtDate"`
	CourtConclusionDate string `json:"courtConclusionDate"`
	Type                string `json:"type"`
	DateEnd             string `json:"dateEnd"`
	CourtName           string `json:"courtName"`
}

type SprVisitN struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type Policy struct {
	PatientId int    `json:"patientId"`
	Company   int    `json:"company"`
	Number    string `json:"number"`
	Series    string `json:"series"`
}
