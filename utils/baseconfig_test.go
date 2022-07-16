package utils

import "testing"

func TestBaseConfig_MaxFloorNum(t *testing.T) {
	bc := &BaseConfig{
		ArticleBusinessID: "",
		XCSRF:             "",
		Cookie:            "",
	}
	num, err := bc.MaxFloorNum("https://bbplanet.bilibili.co/api/planet/comment/commentList")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(num)
}

func TestBaseConfig_ParseURL(t *testing.T) {
	bc := &BaseConfig{
		ArticleBusinessID: "aaaaa",
	}
	url := bc.parseURL("https://localhost:8080/test")
	t.Log(url) // https://localhost:8080/test?articleBusinessId=aaaaa&order=2&pageNum=1000&pageSize=10&scrollId=null
}
