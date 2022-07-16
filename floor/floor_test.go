package floor

import (
	"testing"

	"github.com/demoManito/bilibiliscript/utils"
)

func TestFloor_ParseURL(t *testing.T) {
	f := &Floor{
		Conf: &Conf{
			BaseConfig: utils.BaseConfig{
				URL:               "https://localhost:8080/test",
				ArticleBusinessID: "aaaaa",
			},
		},
	}
	url := f.parseURL(1000)
	t.Log(url) // https://localhost:8080/test?articleBusinessId=aaaaa&order=2&pageNum=1000&pageSize=10&scrollId=null
}
