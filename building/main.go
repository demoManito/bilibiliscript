package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

const tickerDuration = 2 * time.Second

var (
	url = ""

	body = map[string]interface{}{
		"articleBusinessId": "e5339d97ec0e48ec889af1944dd7268e", // TODO: 设置参数
		"atUserList":        make([]interface{}, 0, 0),
		"content":           "1",
	}
)

type RespBody struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func main() {
	ctx := context.Background()
	flag := true
	now := time.Now()
	ticker := time.NewTicker(tickerDuration)
	for {
		select {
		case <-ticker.C:
			switch {
			case flag && now.Add(5*time.Minute).Unix() < time.Now().Unix():
				// 由于接口设置了限流器，所以 5 分钟后请求时长间隔为 5 秒一次
				now = time.Now()
				flag = false
				ticker.Reset(5 * time.Second)
				log.Println("ticker reset 5 s")
			case !flag && now.Add(5*time.Minute).Unix() < time.Now().Unix():
				// 缓冲 5 分钟后，恢复 2 秒一次请求间隔
				now = time.Now()
				flag = true
				ticker.Reset(2 * time.Second)
				log.Println("ticker reset 2 s")
			}

			Building(ctx)
		case <-ctx.Done():
			log.Println("ctx done")
		}
	}
}

func Building(ctx context.Context) {
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		log.Printf("[http request] err: %s \n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF", "csrf-0.5954696332336054")                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   // TODO: 需要设置
	req.Header.Set("Cookie", "buvid3=93EBC383-69C2-9756-1E2C-A1C7FE25433924396infoc; b_nut=1643166524; buvid_fp=e5d0d3c8608036dcdd41f3b32f344f69; buvid4=3AA2FBCC-39C8-73B6-EB62-663C0A70206A73921-022012016-C34HN9S73zrywuG4/UVeSw%3D%3D; mng-go=11ca7ef1fd5a8071d47c837aa36a10c66fae2f2a2adefca684cba067e0cf1543; _AJSESSIONID=67bc6687563f365ba2a9dc0d92fe9f09; username=tangjingyu; b_lsid=63C525C7_181D312239D; b_timer=%7B%22ffp%22%3A%7B%22446.0.fp.risk_93EBC383%22%3A%22181D3122472%22%2C%22333.1019.fp.risk_93EBC383%22%3A%22181D314BCDD%22%7D%7D; SecurityProxySessionID=V1_ZTRhMGZmMDctY2Q5OC0zOTQ0LWI0ZmUtZDljYWYzODNlODBj") // TODO: 需要设置

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Printf("[client do] err: %s \n", err)
	}
	respBody := new(RespBody)
	ioBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(ioBody, respBody)
	if err != nil {
		log.Printf("[resp json unmarshal] err: %s \n", err)
	}
	if respBody.Code != 0 {
		log.Printf("[resp code] message: %s \n", respBody.Message)
	}
}
