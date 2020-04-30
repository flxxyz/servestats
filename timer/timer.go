package timer

import (
	"time"
)

var (
	timer *T
)

type T struct {
	s map[int]*time.Ticker
	i int
}

func (t *T) add(ticker *time.Ticker) int {
	t.s[t.i] = ticker
	return t.i
}

func New(callback func(), interval time.Duration) (i int) {
	i = timer.add(time.NewTicker(interval))
	for range timer.s[i].C {
		callback()
	}
	return
}

func init() {
	timer = &T{
		s: make(map[int]*time.Ticker, 0),
		i: 0,
	}
}
