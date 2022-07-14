package building_test

import (
	"testing"

	"github.com/demoManito/bilibiliscript/building"
)

func TestRun(t *testing.T) {
	building.New("./config.example.yml").Run()
}

func TestRunBuilds(t *testing.T) {
	builds := []*building.Building{
		building.New("./config.example1.yml"), // NOTICE: config 不存在
		building.New("./config.example2.yml"), // NOTICE: config 不存在
	}
	building.RunBuilds(builds)
}
