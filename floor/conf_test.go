package floor

import "testing"

func TestInit(t *testing.T) {
	c := Init("./config.example.yml")
	t.Log(c.URL)
	t.Log(c.FloorNum)
}
