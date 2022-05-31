package types

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
	PatientId   int    `json:"patientId"`
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
	PatientId int    `json:"patientId"`
	DateAdd   string `json:"dateAdd"`
	DockId    int    `json:"dockId"`
	Unit      int    `json:"unit"`
	Zakl      int    `json:"zakl"`
}

type NewRegister struct {
	PatientId  int    `json:"patientId"`
	Reason     string `json:"reason"`
	ExitReason string `json:"exitReason"`
	Section    int    `json:"section"`
	Category   int    `json:"category"`
	Diagnose   string `json:"diagnose"`
	Date       string `json:"date"`
	DockId     int    `json:"dockId"`
}
type NewRegisterTransfer struct {
	PatientId int    `json:"patientId"`
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
