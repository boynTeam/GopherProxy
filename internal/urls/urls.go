package urls

import (
	"github.com/BoynChan/GopherProxy/internal/loadbalance"
)

// Author:Boyn
// Date:2020/8/31

type DynamicUrls struct {
	RegisterAddr string // 注册中心地址
	Strategy     loadbalance.LoadBalance
}

func NewDynamicUrls(urls []string, lbType loadbalance.Type, registerAddr string, param ...string) (*DynamicUrls, error) {
	lb, err := loadbalance.NewStrategy(lbType)
	if err != nil {
		return nil, err
	}
	for _, addr := range urls {
		err := lb.Add(addr)
		if err != nil {
			return nil, err
		}
	}
	return &DynamicUrls{
		RegisterAddr: registerAddr,
		Strategy:     lb,
	}, nil
}

func (s *DynamicUrls) GetNext(key string) (string, error) {
	return s.Strategy.Get(key)
}

func (s *DynamicUrls) GetAllUrl() []string {
	return s.Strategy.GetAll()
}

func (s *DynamicUrls) Watch() {
	blockingChannle := make(chan struct{})
	go func() {
		for {
			// load config()
			<-blockingChannle
		}
	}()
}
