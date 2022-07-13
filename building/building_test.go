package building_test

import (
	"testing"

	"github.com/demoManito/bilibiliscript/building"
)

func TestRun(t *testing.T) {
	building.Run("./config.yml")
}
