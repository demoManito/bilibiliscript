package aliyun_test

import (
	"testing"

	"github.com/demoManito/bilibiliscript/server/aliyun"
)

func TestSendMsg(t *testing.T) {
	sm, _ := aliyun.NewSendMsg("", "", "")
	err := sm.SendMsg("200000")
	if err != nil {
		t.Fatal(err)
	}
}
