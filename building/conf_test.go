package building

import "testing"

func TestInit(t *testing.T) {
	conf := Init("./config.example.yml")
	t.Log(conf.URL)
	t.Log(conf.MaxLimit)
	t.Log(conf.TickerDuration)
	t.Log(conf.TargetFloor)
	t.Log(conf.TargetFloorRule)
}
