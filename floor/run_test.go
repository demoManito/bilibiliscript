package floor_test

import (
	"testing"

	"github.com/demoManito/bilibiliscript/floor"
)

func TestReport(t *testing.T) {
	floor.New("./config.example.yml").Report()
}
