package service

import "sync"

type UserBuyHistory struct {
	history map[int]int
	lock    sync.RWMutex
}

func (p *UserBuyHistory) GetProductBuyCount(productId int) int {
	p.lock.RLock()
	defer p.lock.RUnlock()

	count, _ := p.history[productId]
	return count
}

func (p *UserBuyHistory) Add(productId, count int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	cur, ok := p.history[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}

	p.history[productId] = cur
}
