package building

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	once = &sync.Once{}

	counter int64 = 0
)

type RespBody struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func Run() {
	once.Do(func() {
		conf = Init()
	})

	wait, end, err := await()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Duration(wait))

	ctx := context.Background()
	ticker := time.NewTicker(time.Duration(conf.TickerDuration) * time.Millisecond)
	cleanup := time.NewTicker(10 * time.Second)
loop:
	for {
		select {
		case <-ticker.C:
			if end != 0 && end < time.Now().Unix() {
				break loop
			}
			if atomic.LoadInt64(&counter) > conf.MaxLimit {
				log.Printf("被限流次数超过 %d, 休眠 2 分钟 \n", conf.MaxLimit)
				time.Sleep(2 * time.Minute)
			}
			Building(ctx)
		case <-cleanup.C:
			atomic.StoreInt64(&counter, 0)
		case <-ctx.Done():
			log.Println("ctx done")
		}
	}
	log.Println("👋👋～")
}

func await() (int64, int64, error) {
	if conf.TimingStartTime == "" || conf.TimingEndTime == "" {
		return 0, 0, nil
	}

	now := time.Now()
	start, _ := time.ParseInLocation(conf.TimingStartTime, "", time.Local) // TODO
	end, _ := time.ParseInLocation(conf.TimingEndTime, "", time.Local)     // TODO
	if start.Unix() > now.Unix() {
		wait := now.Unix() - start.Unix()
		return wait, end.Unix(), nil
	}
	if end.Unix() < now.Unix() {
		return 0, 0, errors.New("已结束")
	}
	return 0, end.Unix(), nil
}

func Building(ctx context.Context) {
	body, _ := json.Marshal(map[string]interface{}{
		"articleBusinessId": conf.ArticleBusinessID,
		"atUserList":        make([]interface{}, 0, 0),
		"content":           "1111", // 留言
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, conf.URL, bytes.NewReader(body))
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
		return
	}
	if respBody.Code != 0 {
		log.Printf("[resp code] message: %s \n", respBody.Message)
		switch respBody.Code {
		// TODO
		case 5000:
			atomic.AddInt64(&counter, 1)
		}
		return
	}

	// TODO: 校验是否终止服务
	if conf.TargetFloor == 1000 {

	}

}
