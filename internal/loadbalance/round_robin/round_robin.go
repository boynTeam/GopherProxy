package round_robin

import "sync"

// Author:Boyn
// Date:2020/8/31

type RoundRobinBalance struct {
	curIndex int
	rss      []string
	rw       sync.RWMutex
}

func NewRoundRobin() *RoundRobinBalance {
	return &RoundRobinBalance{}
}

func (r *RoundRobinBalance) Add(s ...string) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.rss = append(r.rss, s...)
	return nil
}

func (r *RoundRobinBalance) Get(s string) (string, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	if len(r.rss) == 0 {
		return "", nil
	}
	lens := len(r.rss) //5
	if r.curIndex >= lens {
		r.curIndex = 0
	}
	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr, nil
}

func (r *RoundRobinBalance) GetAll() []string {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.rss
}

func (r *RoundRobinBalance) Update(newRss []string) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.rss = newRss
	return nil
}
