package building

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncludeFloor(t *testing.T) {
	assert := assert.New(t)

	// test target_floor_rule model 1
	b := &Building{
		Conf: &Conf{
			TargetFloor: &TargetFloor{},
			TargetFloorRule: &TargetFloorRule{
				Enable: true,
				Rule:   1,
				Target: 3,
			},
		},
	}
	assert.False(b.includeFloor(1))
	assert.True(b.includeFloor(3))
	assert.False(b.includeFloor(4))

	// test target_floor_rule model 2
	b = &Building{
		Conf: &Conf{
			TargetFloor: &TargetFloor{},
			TargetFloorRule: &TargetFloorRule{
				Enable: true,
				Rule:   2,
				Target: 3,
			},
		},
	}
	assert.False(b.includeFloor(1111))
	assert.True(b.includeFloor(1113))
	assert.True(b.includeFloor(4131))

	// test target_floor
	b = &Building{
		Conf: &Conf{
			TargetFloor: &TargetFloor{
				Enable: true,
				Nums:   []float64{1, 2},
			},
			TargetFloorRule: &TargetFloorRule{},
		},
	}
	assert.True(b.includeFloor(1))
	assert.False(b.includeFloor(3))
}
