package types

type Patient struct {
	Id         int64  `json:"id"`
	Lname      string `json:"lname"`
	Fname      string `json:"fname"`
	Sname      string `json:"sname"`
	Bday       string `json:"bday"`
	Visibility int    `json:"visibility"`
	Sex        string `json:"sex"`
	Snils      string `json:"snils"`
	Address    string `json:"address"`
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
	Success bool   `json:"success"`
	Error   int    `json:"error"`
	Message string `json:"message"`
}

type NewInvalid struct {
	DoctId       int    `json:"doctId"`
	PatientId    int64  `json:"patientId"`
	DateStart    string `json:"date_start"`
	DateEnd      string `json:"date_end"`
	DateDocument string `json:"date_document"`
	Reason       string `json:"reason"`
	Kind         string `json:"kind"`
	Anomaly      string `json:"anomal"`
	Limit        string `json:"limit"`
	IsInfinity   bool   `json:"isInfinity"`
}

type FindCustody struct {
	PatientId int64  `json:"patientId"`
	DoctId    int    `json:"doctId"`
	DateStart string `json:"dateStart"`
	DateEnd   string `json:"dateEnd"`
	Who       string `json:"who"`
}
