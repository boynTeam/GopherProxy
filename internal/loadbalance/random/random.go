package random

import (
	"math/rand"
	"sync"
)

// Author:Boyn
// Date:2020/8/31

type RandomBalance struct {
	rss []string
	rw  sync.RWMutex
}

func NewRandomBalance() *RandomBalance {
	return &RandomBalance{}
}

func (r *RandomBalance) Add(s ...string) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.rss = append(r.rss, s...)
	return nil
}

func (r *RandomBalance) Get(s string) (string, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	if len(r.rss) == 0 {
		return "", nil
	}
	randomKey := rand.Intn(len(r.rss))
	return r.rss[randomKey], nil
}

func (r *RandomBalance) GetAll() []string {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.rss
}

func (r *RandomBalance) Update(newRss []string) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.rss = newRss
	return nil
}
