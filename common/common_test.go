package common

import "testing"

func TestRandString(t *testing.T) {
	s := []int{
		1,
		2,
		3,
		4,
		5,
	}

	for _, v := range s {
		rs := RandString(v)

		if len(rs) != v {
			t.FailNow()
		}
	}
}
