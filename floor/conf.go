package floor

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/demoManito/bilibiliscript/utils"
)

type Conf struct {
	utils.BaseConfig `yaml:",inline"`

	FloorNum int64 `yaml:"floor_num"`
}

func Init(filename string) *Conf {
	file, err := os.ReadFile(filename)
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
