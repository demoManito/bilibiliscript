package building

import "testing"

func TestInit(t *testing.T) {
	conf := Init("./config.yml")
	t.Log(conf.MaxLimit)
	t.Log(conf.TickerDuration)
}
