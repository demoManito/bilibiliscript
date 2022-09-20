package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

// BaseConfig 基础配置字段
type BaseConfig struct {
	URL               string `yaml:"url"` // 对应接口
	ArticleBusinessID string `yaml:"article_business_id"`
	XCSRF             string `yaml:"xcsrf"`
	Cookie            string `yaml:"cookie"`
}

func (bc *BaseConfig) SetReqHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-CSRF", bc.XCSRF)
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("sec-ch-ua", "\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"103\", \"Chromium\";v=\"103\"")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.Header.Set("X-AppKey", "ops.teamwork.portal")
	req.Header.Set("Cookie", bc.Cookie)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-UserType", "1")
}

// MaxFloorNum 当前最大楼层
func (bc *BaseConfig) MaxFloorNum(url string) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, bc.parseURL(url), nil)
	if err != nil {
		log.Printf("[request err] err: %s", err)
		return 0, err
	}
	bc.SetReqHeader(req)

	resp := new(Resp)
	response, err := new(http.Client).Do(req)
	if err != nil {
		log.Printf("[http err] http client do err: %s", err)
		return 0, err
	}
	ioBody, _ := io.ReadAll(EncodingBody(response))
	err = json.Unmarshal(ioBody, &resp)
	if err != nil {
		log.Printf("[unmarshal err] resp json unmarshal err: %s", err)
		return 0, err
	}
	if resp.Code != 0 {
		log.Printf("[resp err] code: %d; message: %s; err: %s \n", resp.Code, resp.Message, err)
		return 0, err
	}

	list := resp.Data["commentReplyList"]
	jsonList, err := json.Marshal(list)
	if err != nil {
		log.Printf("[marshal err] comment_reply_list json marshal err: %s", err)
		return 0, err
	}
	floors := make([]FloorInfo, 0)
	err = json.Unmarshal(jsonList, &floors)
	if err != nil {
		log.Printf("[unmarshal err] floors json unmarshal err: %s", err)
		return 0, err
	}

	if len(floors) <= 0 {
		return 0, nil
	}
	return floors[0].FloorNum, nil
}

func (bc *BaseConfig) parseURL(u string) string {
	up, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	q := up.Query()
	q.Set("articleBusinessId", bc.ArticleBusinessID)
	q.Set("pageSize", "10")
	q.Set("pageNum", "1")
	q.Set("order", "1")
	q.Set("scrollId", "null")
	up.RawQuery = q.Encode()
	return up.String()
}
