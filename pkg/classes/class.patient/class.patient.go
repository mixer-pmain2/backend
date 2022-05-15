package class_patient

import "pmain2/internal/models"

type PatientClass struct {
	Personal models.Patient
	LastReg  models.FindUchetS
}

func Init() {

}

func (p *PatientClass) Get(id int) {

}
