package building

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var conf *Conf

type Conf struct {
	URL               string `yaml:"url"`
	ArticleBusinessID string `yaml:"article_business_id"`
	XCSRF             string `yaml:"xcsrf"`
	Cookie            string `yaml:"cookie"`
	TargetFloor       int64  `yaml:"target_floor,omitempty"`      // 目标楼层, 盖中此楼将退出盖楼
	TimingStartTime   string `yaml:"timing_start_time,omitempty"` // 定时开始时间
	TimingEndTime     string `yaml:"timing_end_time,omitempty"`   // 定时结束时间
}

func Init() *Conf {
	file, err := ioutil.ReadFile("./config.yml")
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
