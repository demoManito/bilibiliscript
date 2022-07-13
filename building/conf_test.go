package building

import "testing"

func TestInit(t *testing.T) {
	conf := Init()
	t.Log(conf.MaxLimit)
	t.Log(conf.TickerDuration)
}
