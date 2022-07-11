package report

var (
	statusType = struct {
		new      string
		progress string
		done     string
		error    string
	}{
		"NEW",
		"PROGRESS",
		"DONE",
		"ERROR",
	}

	reportType = struct {
		ReceptionLog           string
		VisitsPerPeriod        string
		AdmittedToTheHospital  string
		DischargedFromHospital string
		Unvisited              string
		Registered             string
		Deregistered           string
		ConsistingOnTheSite    string
	}{
		ReceptionLog:           "ReceptionLog",
		VisitsPerPeriod:        "VisitsPerPeriod",
		AdmittedToTheHospital:  "AdmittedToTheHospital",
		DischargedFromHospital: "DischargedFromHospital",
		Unvisited:              "Unvisited",
		Registered:             "Registered",
		Deregistered:           "Deregistered",
		ConsistingOnTheSite:    "ConsistingOnTheSite",
	}
)

type orderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   int    `json:"error"`
}

type reportFilter struct {
	RangeDate    []string `json:"rangeDate"`
	DateStart    string   `json:"dateStart"`
	DateEnd      string   `json:"dateEnd"`
	Category     int      `json:"category"`
	TypeCategory string   `json:"typeCategory"`
	RangeSection []int    `json:"rangeSection"`
}

type reportParams struct {
	Id      int          `json:"id"`
	Code    string       `json:"code"`
	UserId  int          `json:"userId"`
	Unit    int          `json:"unit"`
	Filters reportFilter `json:"filters"`
	Status  string       `json:"status"`
	Date    string       `json:"date"`
}

type reportData struct {
	Title string
	Data  []interface{}
}
