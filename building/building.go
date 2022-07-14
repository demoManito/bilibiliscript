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

type Building struct {
	Conf *Conf

	ctx      context.Context
	counter  int64
	floorNum int64
	done     chan struct{}
}

func New(fileName string) *Building {
	return &Building{
		Conf:    Init(fileName),
		counter: 0,
		done:    make(chan struct{}, 0),
		ctx:     context.Background(),
	}
}

func (b *Building) Run() {
	wait, end := b.await()
	log.Printf("距开始盖楼需等待：%d 秒 \n", wait)
	time.Sleep(time.Duration(wait) * time.Second)

	ticker := time.NewTicker(time.Duration(b.Conf.TickerDuration) * time.Millisecond)
	cleanup := time.NewTicker(10 * time.Second)
loop:
	for {
		select {
		case <-ticker.C:
			if end != 0 && end < time.Now().Unix() {
				break loop
			}
			if b.Conf.TickerDuration >= 1000 && atomic.LoadInt64(&b.counter) > b.Conf.MaxLimit {
				log.Printf("被限流次数超过 %d 次, 休眠 11 秒 \n", b.Conf.MaxLimit)
				atomic.StoreInt64(&b.counter, 0)
				time.Sleep(11 * time.Second)
			}
			go b.building(b.ctx)
		case <-cleanup.C:
			atomic.StoreInt64(&b.counter, 0)
		case <-b.done:
			log.Println("done")
			break loop
		case <-b.ctx.Done():
			log.Println("ctx done")
			break loop
		}
	}
	log.Printf("盖了 %d 楼, 👋👋～", atomic.LoadInt64(&b.floorNum))
}

func (b *Building) await() (int64, int64) {
	if b.Conf.TimingStartTime == "" || b.Conf.TimingEndTime == "" {
		return 0, 0
	}
	log.Printf("hi～ start_time: %s, end_time: %s \n", b.Conf.TimingStartTime, b.Conf.TimingEndTime)

	now := time.Now()
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", b.Conf.TimingStartTime, time.Local)
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", b.Conf.TimingEndTime, time.Local)
	if start.Unix() > now.Unix() {
		wait := start.Unix() - now.Unix()
		return wait, end.Unix()
	}
	if end.Unix() < now.Unix() {
		log.Fatal(errors.New("已结束"))
	}
	return 0, end.Unix()
}

func (b *Building) building(ctx context.Context) {
	body, _ := json.Marshal(map[string]interface{}{
		"articleBusinessId": b.Conf.ArticleBusinessID,
		"atUserList":        make([]interface{}, 0, 0),
		"content":           rand.Int31n(10), // 留言
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.Conf.URL, bytes.NewReader(body))
	if err != nil {
		log.Printf("[http request] err: %s \n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF", b.Conf.XCSRF)
	req.Header.Set("Cookie", b.Conf.Cookie)

	var resp struct {
		Code    int64                  `json:"code"`
		Data    map[string]interface{} `json:"data"`
		Message string                 `json:"message"`
	}
	response, err := new(http.Client).Do(req)
	if err != nil {
		log.Printf("[client do] err: %s \n", err)
	}
	ioBody, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(ioBody, &resp)
	if err != nil {
		log.Printf("[resp json unmarshal] err: %s \n", err)
		return
	}
	if resp.Code != 0 {
		log.Printf("[resp err] code: %d; message: %s \n", resp.Code, resp.Message)
		switch resp.Code {
		case 90005: // 太快啦~不要刷啦~
			atomic.AddInt64(&b.counter, 1)
		}
		return
	}
	atomic.AddInt64(&b.floorNum, 1)
	// 盖中目标楼层，终止盖楼
	if floorNum, ok := resp.Data["floorNum"]; ok && b.includeFloor(floorNum.(float64)) {
		b.done <- struct{}{}
		log.Printf("恭喜🎉🎉🎉～ %.0f层盖中啦～ \n", b.Conf.TargetFloor)
	}
}

func (b *Building) includeFloor(floorNum float64) bool {
	if len(b.Conf.TargetFloor) != 0 {
		for _, tf := range b.Conf.TargetFloor {
			if tf == floorNum {
				return true
			}
		}
	}
	if len(b.Conf.TargetFloorRule) != 0 {
		switch b.Conf.TargetFloorRule["rule"] {
		case targetFloorRuleMOD:
			if int(floorNum)%b.Conf.TargetFloorRule["target"] == 0 {
				return true
			}
		}
	}
	return false
}

// RunBuilds 同时盖多个贴的楼
func RunBuilds(builds []*Building) {
	var wg sync.WaitGroup
	for _, build := range builds {
		wg.Add(1)
		go func(b *Building) {
			defer wg.Done()
			b.Run()
		}(build)
	}
	wg.Wait()
}
