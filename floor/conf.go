package floor

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/demoManito/bilibiliscript/utils"
)

type Conf struct {
	utils.ScriptConfig `yaml:",inline"`

	FloorNum int64 `yaml:"floor_num"`
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
