package client

import (
	"io/ioutil"

	"github.com/1uvu/fabric-sdk-client/types"
	"gopkg.in/yaml.v2"
)

//
// 读入 *.yaml 配置文件
//

type clientConfig struct {
	Admin adminConfig `yaml:"admin"`
	App   appConfig   `yaml:"app"`
}

type adminConfig struct {
	Params types.AdminParams `yaml:"params"`
	Envs   []types.EnvPair   `yaml:"envs,flow"`
}

type appConfig struct {
	Params types.AppParams `yaml:"params"`
	Envs   []types.EnvPair `yaml:"envs,flow"`
}

func newClientConfig(confPath string) (*clientConfig, error) {

	conf, err := getClientConfig(confPath)

	if err != nil {
		return nil, err
	}

	return conf, nil
}

func getClientConfig(confPath string) (*clientConfig, error) {
	conf := new(clientConfig)
	confFile, err := ioutil.ReadFile(confPath)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(confFile, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
