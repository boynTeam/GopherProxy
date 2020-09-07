package loadbalance

import (
	"errors"
)

// Author:Boyn
// Date:2020/8/31

type Type int

const (
	ConsistentHash = iota + 1
	Random
	RoundRobin
)

type LoadBalance interface {
	Add(...ConfigValue) error
	Get(string) (string, error)
	GetAll() []string
	Update([]ConfigValue) error
}

func NewStrategy(lbType Type) (LoadBalance, error) {
	switch lbType {
	case Random:
		return NewRandomBalance(), nil
	case RoundRobin:
		return NewRoundRobin(), nil
	case ConsistentHash:
		return NewConsistentHash(), nil
	default:
		return nil, errors.New("unsupported type")
	}
}

func NewStrategyWithConf(lbType Type, mConf Config) (LoadBalance, error) {
	configValues, err := mConf.GetConf()
	if err != nil {
		return nil, err
	}
	var b LoadBalance
	switch lbType {
	case Random:
		b = NewRandomBalance()
	case RoundRobin:
		b = NewRoundRobin()
	case ConsistentHash:
		b = NewConsistentHash()
	default:
		return nil, errors.New("unsupported type")
	}
	mConf.Attach(b)
	b.Add(configValues...)
	return b, nil
}
