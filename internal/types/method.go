package types

func (v *NewVisit) Normalize() {
	if v.Unit == 1 {
		v.Unit = 0
	}
	if v.Home {
		v.Unit += 1
	}
}
