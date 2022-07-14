package building

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const (
	targetFloorRuleMOD = iota + 1 // 目标数的倍数
)

type Conf struct {
	MaxLimit          int64          `yaml:"max_limit"`       // 接口提示盖楼频繁最大限制次数，达到限制将休眠 x 秒后请求
	TickerDuration    int64          `yaml:"ticker_duration"` // 间隔多少毫秒轮询一次
	URL               string         `yaml:"url"`
	ArticleBusinessID string         `yaml:"article_business_id"`
	XCSRF             string         `yaml:"xcsrf"`
	Cookie            string         `yaml:"cookie"`
	TargetFloor       []float64      `yaml:"target_floor,omitempty"`      // 目标楼层, 盖中此楼将退出盖楼
	TargetFloorRule   map[string]int `yaml:"target_floor_rule,omitempty"` // 带规则匹配目标楼层
	TimingStartTime   string         `yaml:"timing_start_time,omitempty"` // 定时开始时间
	TimingEndTime     string         `yaml:"timing_end_time,omitempty"`   // 定时结束时间
}

func Init(filename string) *Conf {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	conf := new(Conf)
	err = yaml.Unmarshal(file, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
