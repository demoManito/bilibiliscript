package building

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/demoManito/bilibiliscript/utils"
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
	b.waiter()

	ticker := time.NewTicker(time.Duration(b.Conf.TickerDuration) * time.Millisecond)
	cleanup := time.NewTicker(5 * time.Second)
loop:
	for {
		select {
		case <-ticker.C:
			go b.building(b.ctx)
		case <-cleanup.C:
			b.cleanup()
		case <-b.done:
			b.clear()
			break loop
		case <-b.ctx.Done():
			b.clear()
			break loop
		}
	}
	log.Printf("成功盖了 %d 层 👋👋～", atomic.LoadInt64(&b.floorNum))
}

func (b *Building) waiter() {
	switch {
	case b.Conf.TriggerBuilding.Enable:
		b.triggerFloor(b.Conf.TriggerBuilding.URL, b.Conf.TriggerBuilding.Num)
	case b.Conf.TargetFloorScope.Enable:
		b.triggerFloor(b.Conf.TargetFloorScope.URL, b.Conf.TargetFloorScope.MIN)
	case b.Conf.Timing.Enable:
		b.timing()
	}

	log.Println("🏠 开始盖楼啦～")
}

func (b *Building) triggerFloor(url string, num int64) {
	log.Printf("正在等待楼层 %d 生成...", num)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if b.isTriggerBuilding(url, num) {
				ticker.Stop()
				return
			}
		}
	}
}

func (b *Building) timing() {
	now := time.Now()
	if b.Conf.Timing.StartTime != "" {
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", b.Conf.Timing.StartTime, time.Local)
		if start.Unix() > now.Unix() {
			wait := start.Unix() - now.Unix()
			log.Printf("距开始盖楼需等待：%d 秒 \n", wait)
			<-time.NewTimer(time.Duration(wait) * time.Second).C
		}
	}
}

func (b *Building) building(ctx context.Context) {
	if b.Conf.TickerDuration >= 1000 && atomic.LoadInt64(&b.counter) > b.Conf.MaxLimit {
		log.Printf("被限流次数超过 %d 次, 休眠 11 秒 \n", b.Conf.MaxLimit)
		atomic.StoreInt64(&b.counter, 0)
		<-time.NewTimer(9 * time.Second).C
	}

	body, _ := json.Marshal(map[string]interface{}{
		"articleBusinessId": b.Conf.ArticleBusinessID,
		"atUserList":        make([]interface{}, 0, 0),
		"content":           rand.Int31n(10), // 留言
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.Conf.URL, bytes.NewReader(body))
	if err != nil {
		log.Printf("[request err] err: %s \n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF", b.Conf.XCSRF)
	req.Header.Set("Cookie", b.Conf.Cookie)

	resp := new(utils.Resp)
	response, err := new(http.Client).Do(req)
	if err != nil {
		log.Printf("[http err] client do err: %s \n", err)
		return
	}
	ioBody, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(ioBody, &resp)
	if err != nil {
		log.Printf("[unmarshal err] resp json unmarshal err: %s \n", err)
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
		log.Printf("恭喜🎉🎉🎉 %.0f 层盖中啦～ \n", floorNum)
	}
}

func (b *Building) includeFloor(floorNum float64) bool {
	if b.Conf.TargetFloor.Enable {
		for _, tf := range b.Conf.TargetFloor.Nums {
			if tf == floorNum {
				return true
			}
		}
	}
	if b.Conf.TargetFloorRule.Enable {
		target := b.Conf.TargetFloorRule.Target
		switch b.Conf.TargetFloorRule.Rule {
		case targetFloorRuleMOD:
			if int(floorNum)%target == 0 {
				return true
			}
		case targetFloorRuleInclude:
			if strings.Contains(strconv.FormatFloat(floorNum, 'f', 2, 64), strconv.FormatInt(int64(target), 10)) {
				return true
			}
		}
	}
	if b.Conf.TargetFloorScope.Enable {
		if int(floorNum) >= int(b.Conf.TargetFloorScope.MAX) {
			return true
		}
	}
	return false
}

func (b *Building) cleanup() {
	atomic.StoreInt64(&b.counter, 0)

	// 扫描定时任务结束时间
	if b.Conf.Timing.Enable && b.Conf.Timing.EndTime != "" {
		go func() {
			end, _ := time.ParseInLocation("2006-01-02 15:04:05", b.Conf.Timing.EndTime, time.Local)
			if end.Unix() <= time.Now().Unix() {
				log.Println("⏰ 盖楼结束啦！")
				b.done <- struct{}{}
			}
		}()
	}
}

func (b *Building) clear() {
	close(b.done)
	log.Println("下次再见～ 🐛🐛🐛")
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
