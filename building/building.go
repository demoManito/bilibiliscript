package building

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	once = &sync.Once{}

	done = make(chan struct{}, 0)

	counter int64 = 0
)

func Run(file string) {
	once.Do(func() {
		conf = Init(file)
	})

	wait, end, err := await()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("距开始盖楼需等待：%d 秒", wait)
	time.Sleep(time.Duration(wait) * time.Second)

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
				log.Printf("被限流次数超过 %d 次, 休眠 10 秒 \n", conf.MaxLimit)
				atomic.StoreInt64(&counter, 0)
				time.Sleep(10 * time.Second)
			}
			go building(ctx)
		case <-cleanup.C:
			log.Println("reset counter")
			atomic.StoreInt64(&counter, 0)
		case <-done:
			log.Println("done")
			break loop
		case <-ctx.Done():
			log.Println("ctx done")
			break loop
		}
	}
	log.Println("👋👋～")
}

func await() (int64, int64, error) {
	if conf.TimingStartTime == "" || conf.TimingEndTime == "" {
		return 0, 0, nil
	}
	log.Printf("hi～ start_time: %s, end_time: %s \n", conf.TimingStartTime, conf.TimingEndTime)

	now := time.Now()
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", conf.TimingStartTime, time.Local)
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", conf.TimingEndTime, time.Local)
	if start.Unix() > now.Unix() {
		wait := start.Unix() - now.Unix()
		return wait, end.Unix(), nil
	}
	if end.Unix() < now.Unix() {
		return 0, 0, errors.New("已结束")
	}
	return 0, end.Unix(), nil
}

type RespBody struct {
	Code    int64                  `json:"code"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}

func building(ctx context.Context) {
	body, _ := json.Marshal(map[string]interface{}{
		"articleBusinessId": conf.ArticleBusinessID,
		"atUserList":        make([]interface{}, 0, 0),
		"content":           rand.Int31n(10), // 留言
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, conf.URL, bytes.NewReader(body))
	if err != nil {
		log.Printf("[http request] err: %s \n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF", conf.XCSRF)
	req.Header.Set("Cookie", conf.Cookie)

	response, err := new(http.Client).Do(req)
	if err != nil {
		log.Printf("[client do] err: %s \n", err)
	}
	resp := new(RespBody)
	ioBody, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(ioBody, resp)
	if err != nil {
		log.Printf("[resp json unmarshal] err: %s \n", err)
		return
	}
	if resp.Code != 0 {
		log.Printf("[resp err] code: %d; message: %s \n", resp.Code, resp.Message)
		switch resp.Code {
		case 90005: // 太快啦~不要刷啦~
			atomic.AddInt64(&counter, 1)
		}
		return
	}
	// 盖中目标楼层，终止盖楼
	floorNum, ok := resp.Data["floorNum"]
	if ok {
		if conf.TargetFloor != 0 && conf.TargetFloor == floorNum.(float64) {
			done <- struct{}{}
			log.Printf("恭喜🎉🎉🎉～ %.0f层盖中啦～ \n", conf.TargetFloor)
		}
	}
}
