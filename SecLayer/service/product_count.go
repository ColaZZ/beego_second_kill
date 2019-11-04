package service

import "sync"

type ProductCountMgr struct {
	ProductCount map[int]int
	Lock         sync.RWMutex
}

func NewProductCountMgr() (productMgr *ProductCountMgr) {
	productMgr = &ProductCountMgr{
		ProductCount: make(map[int]int, 128),
	}
	return
}

func (p *ProductCountMgr) Count(productId int) (count int) {
	p.Lock.RLock()
	defer p.Lock.RUnlock()

	count, _ = p.ProductCount[productId]
	return
}

func (p *ProductCountMgr) Add(productId, count int) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	cur, ok := p.ProductCount[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	p.ProductCount[productId] = cur
}
