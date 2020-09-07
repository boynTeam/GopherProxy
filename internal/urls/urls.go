package urls

import (
	"fmt"

	"github.com/BoynChan/GopherProxy/internal/loadbalance"
	"github.com/spf13/viper"
)

// Author:Boyn
// Date:2020/8/31

type DynamicUrls struct {
	RegisterAddr string // 注册中心地址
	ServiceName  string
	Strategy     loadbalance.LoadBalance
	Config       loadbalance.Config
}

func NewDynamicUrls(serviceName string, lbType loadbalance.Type, registerAddr string, param ...string) (*DynamicUrls, error) {
	httpPrefix := viper.GetString("Zk.HttpPrefix")
	conf, err := loadbalance.NewZkConf("%s", fmt.Sprintf("%s/%s", httpPrefix, serviceName), []string{registerAddr})
	if err != nil {
		return nil, err
	}
	lb, err := loadbalance.NewStrategyWithConf(lbType, conf)
	if err != nil {
		return nil, err
	}
	return &DynamicUrls{
		RegisterAddr: registerAddr,
		Config:       conf,
		Strategy:     lb,
	}, nil
}

func (s *DynamicUrls) GetNext(key string) (string, error) {
	return s.Strategy.Get(key)
}

func (s *DynamicUrls) GetAllUrl() []string {
	return s.Strategy.GetAll()
}
