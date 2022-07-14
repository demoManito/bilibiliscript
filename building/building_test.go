package building

import "testing"

func TestIncludeFloor(t *testing.T) {
	b := &Building{
		Conf: &Conf{
			TargetFloorRule: map[string]int{"rule": 1, "target": 3},
		},
	}
	t.Log(b.includeFloor(1)) // false
	t.Log(b.includeFloor(3)) // true
	t.Log(b.includeFloor(4)) // false

	b = &Building{
		Conf: &Conf{
			TargetFloor: []float64{1, 2},
		},
	}
	t.Log(b.includeFloor(1)) // true
	t.Log(b.includeFloor(3)) // false
}
