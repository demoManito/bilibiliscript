package building

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/demoManito/bilibiliscript/utils"
)

func (b *Building) triggerBuilding() bool {
	req, err := http.NewRequest(http.MethodGet, b.parseURL(), nil)
	if err != nil {
		log.Printf("[request err] err: %s", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF", b.Conf.XCSRF)
	req.Header.Set("Cookie", b.Conf.Cookie)

	resp := new(utils.Resp)
	response, err := new(http.Client).Do(req)
	if err != nil {
		log.Printf("[http err] http client do err: %s", err)
		return false
	}
	ioBody, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(ioBody, &resp)
	if err != nil {
		log.Printf("[unmarshal err] resp json unmarshal err: %s", err)
		return false
	}
	if resp.Code != 0 {
		log.Printf("[resp err] code: %d; message: %s; err: %s \n", resp.Code, resp.Message, err)
		return false
	}

	list := resp.Data["commentReplyList"]
	jsonList, err := json.Marshal(list)
	if err != nil {
		log.Printf("[marshal err] comment_reply_list json marshal err: %s", err)
		return false
	}
	floors := make([]*utils.FloorInfo, 0)
	err = json.Unmarshal(jsonList, &floors)
	if err != nil {
		log.Printf("[unmarshal err] floors json unmarshal err: %s", err)
		return false
	}

	if b.Conf.TriggerBuilding.Num <= floors[0].FloorNum {
		return true
	}
	return false
}

func (b *Building) parseURL() string {
	up, err := url.Parse(b.Conf.TriggerBuilding.URL)
	if err != nil {
		log.Fatal(err)
	}
	q := up.Query()
	q.Set("articleBusinessId", b.Conf.ArticleBusinessID)
	q.Set("pageSize", "10")
	q.Set("pageNum", "1")
	q.Set("order", "1")
	q.Set("scrollId", "null")
	up.RawQuery = q.Encode()
	return up.String()
}
