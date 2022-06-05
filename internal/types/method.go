package types

import (
	"strings"
	"time"

	"pmain2/pkg/utils"
)

func (v *NewVisit) Normalize() {
	if v.Unit == 1 {
		v.Unit = 0
	}
	if v.Home {
		v.Unit += 1
	}
}

func (v *NewProf) Normalize() {
	if v.Unit == 1 {
		v.Unit = 0
	}
	if v.Home {
		v.Unit += 1
	}
}

func (m *Patient) Serialize() error {
	var err error
	m.Lname, err = utils.ToUTF8(m.Lname)
	if err != err {
		return err
	}
	m.Fname, err = utils.ToUTF8(m.Fname)
	if err != err {
		return err
	}
	m.Sname, err = utils.ToUTF8(m.Sname)
	if err != err {
		return err
	}
	m.Sex, err = utils.ToUTF8(m.Sex)
	if err != err {
		return err
	}
	m.Address, err = utils.ToUTF8(m.Address)
	if err != err {
		return err
	}
	m.Lname = strings.ReplaceAll(m.Lname, " ", "")
	m.Fname = strings.ReplaceAll(m.Fname, " ", "")
	m.Sname = strings.ReplaceAll(m.Sname, " ", "")
	m.Snils = strings.ReplaceAll(m.Snils, " ", "")
	m.Address = strings.ReplaceAll(m.Address, " ", "")
	bdayTime, err := time.Parse(time.RFC3339, m.Bday)
	if err != err {
		return err
	}
	m.Bday = utils.ToDate(bdayTime)
	return nil
}
