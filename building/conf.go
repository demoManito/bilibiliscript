package building

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/demoManito/bilibiliscript/utils"
)

const (
	targetFloorRuleMOD     = iota + 1 // 目标数的倍数
	targetFloorRuleInclude            // 包含某个数
)

// Timing 定时相关
type Timing struct {
	Enable    bool   `yaml:"enable"`     // 是否开启
	StartTime string `yaml:"start_time"` // 定时开始时间
	EndTime   string `yaml:"end_time"`   // 定时结束时间
}

type Conf struct {
	utils.ScriptConfig `yaml:",inline"`

	MaxLimit        int64          `yaml:"max_limit"`                   // 接口提示盖楼频繁最大限制次数，达到限制将休眠 x 秒后请求
	TickerDuration  int64          `yaml:"ticker_duration"`             // 间隔多少毫秒轮询一次
	Timing          *Timing        `yaml:"timing"`                      // 定时
	TargetFloor     []float64      `yaml:"target_floor,omitempty"`      // 目标楼层, 盖中此楼将退出盖楼
	TargetFloorRule map[string]int `yaml:"target_floor_rule,omitempty"` // 带规则匹配目标楼层
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
