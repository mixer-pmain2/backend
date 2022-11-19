package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Version string

func (v Version) split() []int {
	res := []int{0, 0, 0}
	for i, val := range strings.Split(fmt.Sprintf("%s", v), ".") {
		res[i], _ = strconv.Atoi(val)
	}
	return res
}

func (v Version) IsValid() (bool, error) {
	sp := strings.Split(fmt.Sprintf("%s", v), ".")
	for _, s := range sp {
		_, err := strconv.Atoi(s)
		if err != nil {
			return false, fmt.Errorf("version is not valid code, %s", v)
		}
	}
	if len(sp) != 3 {
		return false, fmt.Errorf("version is not valid code, %s", v)
	}

	return true, nil
}

func (v Version) IsHigh(s Version) (bool, error) {
	if _, err := v.IsValid(); err != nil {
		return false, err
	}
	if _, err := s.IsValid(); err != nil {
		return false, err
	}

	//версия "v" старше "s"
	for i := range v.split() {
		if s.split()[i] < v.split()[i] {
			return true, nil
		}
	}

	return false, nil
}
