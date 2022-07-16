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

// TriggerFloor 楼层触发器
type TriggerFloor struct {
	Enable bool   `yaml:"enable"` // 是否开启
	URL    string `yaml:"url"`    // 评论列表接口
	Num    int64  `yaml:"num"`    // 触发开始盖楼的楼层数
}

// TargetFloor 盖中目标楼层
type TargetFloor struct {
	Enable bool      `yaml:"enable"`
	Nums   []float64 `yaml:"nums"` // 目标楼层, 盖中任意一楼层将退出盖楼
}

// TargetFloorRule 带规则匹配目标楼层
type TargetFloorRule struct {
	Enable bool `yaml:"enable"`
	Rule   int  `yaml:"rule"`
	Target int  `yaml:"target"`
}

// TargetFloorScope 盖区间内的楼
type TargetFloorScope struct {
	Enable bool   `yaml:"enable"`
	URL    string `yaml:"url"`
	MIN    int64  `yaml:"min"`
	MAX    int64  `yaml:"max"`
}

type Conf struct {
	utils.BaseConfig `yaml:",inline"`

	MaxLimit       int64 `yaml:"max_limit"`       // 接口提示盖楼频繁最大限制次数，达到限制将休眠 x 秒后请求
	TickerDuration int64 `yaml:"ticker_duration"` // 间隔多少毫秒轮询一次

	Timing           *Timing           `yaml:"timing"`             // 定时器
	TriggerBuilding  *TriggerFloor     `yaml:"trigger_building"`   // 楼层触发器
	TargetFloor      *TargetFloor      `yaml:"target_floor"`       // 盖中指定目标楼层
	TargetFloorRule  *TargetFloorRule  `yaml:"target_floor_rule"`  // 盖中规则匹配上的目标楼层
	TargetFloorScope *TargetFloorScope `yaml:"target_floor_scope"` // 盖中目标楼层返回
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
