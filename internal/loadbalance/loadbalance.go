package loadbalance

import (
	"errors"

	"github.com/BoynChan/GopherProxy/internal/loadbalance/consistent_hash"
	"github.com/BoynChan/GopherProxy/internal/loadbalance/random"
	"github.com/BoynChan/GopherProxy/internal/loadbalance/round_robin"
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
	Add(...string) error
	Get(string) (string, error)
	GetAll() []string
	Update([]string) error
}

func NewStrategy(lbType Type) (LoadBalance, error) {
	switch lbType {
	case Random:
		return random.NewRandomBalance(), nil
	case RoundRobin:
		return round_robin.NewRoundRobin(), nil
	case ConsistentHash:
		return consistent_hash.NewConsistentHash(), nil
	default:
		return nil, errors.New("unsupported type")
	}
}
