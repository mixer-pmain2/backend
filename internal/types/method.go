package types

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
