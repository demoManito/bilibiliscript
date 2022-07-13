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
		"articleBusinessId": "", // TODO: 设置参数
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
	req.Header.Set("X-CSRF", "") // TODO: 需要设置
	req.Header.Set("Cookie", "") // TODO: 需要设置

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
