package building

import "testing"

func TestIncludeFloor(t *testing.T) {
	t.Log("model 1")
	b := &Building{
		Conf: &Conf{
			TargetFloorRule: map[string]int{"rule": 1, "target": 3},
		},
	}
	t.Log(b.includeFloor(1)) // false
	t.Log(b.includeFloor(3)) // true
	t.Log(b.includeFloor(4)) // false

	t.Log("model 2")
	b = &Building{
		Conf: &Conf{
			TargetFloorRule: map[string]int{"rule": 2, "target": 3},
		},
	}
	t.Log(b.includeFloor(1111)) // false
	t.Log(b.includeFloor(1113)) // true
	t.Log(b.includeFloor(4131)) // true

	t.Log("model array")
	b = &Building{
		Conf: &Conf{
			TargetFloor: []float64{1, 2},
		},
	}
	t.Log(b.includeFloor(1)) // true
	t.Log(b.includeFloor(3)) // false
}
