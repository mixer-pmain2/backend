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
		ReceptionLog                   string
		VisitsPerPeriod                string
		AdmittedToTheHospital          string
		DischargedFromHospital         string
		Unvisited                      string
		Registered                     string
		Deregistered                   string
		ConsistingOnTheSite            string
		ThoseInTheHospital             string
		HospitalTreatment              string
		AmbulatoryTreatment            string
		PBSTIN                         string
		TakenOnADN                     string
		TakenFromADN                   string
		TakenForADNAccordingToClinical string
		ProtocolUKL                    string
	}{
		ReceptionLog:                   "ReceptionLog",
		VisitsPerPeriod:                "VisitsPerPeriod",
		AdmittedToTheHospital:          "AdmittedToTheHospital",
		DischargedFromHospital:         "DischargedFromHospital",
		Unvisited:                      "Unvisited",
		Registered:                     "Registered",
		Deregistered:                   "Deregistered",
		ConsistingOnTheSite:            "ConsistingOnTheSite",
		ThoseInTheHospital:             "ThoseInTheHospital",
		HospitalTreatment:              "HospitalTreatment",
		AmbulatoryTreatment:            "AmbulatoryTreatment",
		PBSTIN:                         "PBSTIN",
		TakenOnADN:                     "TakenOnADN",
		TakenFromADN:                   "TakenFromADN",
		TakenForADNAccordingToClinical: "TakenForADNAccordingToClinical",
		ProtocolUKL:                    "ProtocolUKL",
	}
)

type orderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   int    `json:"error"`
}

type reportFilter struct {
	Id           int      `json:"id"`
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
