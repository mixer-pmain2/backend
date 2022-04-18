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
	DockID      int    `json:"dockID"`
	SRC         int    `json:"src"`
}

type NewSRC struct {
	PatientId int    `json:"patientId"`
	DateAdd   string `json:"dateAdd"`
	DockId    int    `json:"dockId"`
	Unit      int    `json:"unit"`
	Zakl      int    `json:"zakl"`
}

type HttpResponse struct {
	Success bool   `json:"success"`
	Error   int    `json:"error"`
	Message string `json:"message"`
}
