package timer

import (
	"time"
)

var (
	queue chan *T
)

type T struct {
	Interval time.Duration
	Callback func()
}

func Run() {
	for {
		select {
		case t, ok := <-queue:
			if ok {
				ticker := time.NewTicker(t.Interval)
				for range ticker.C {
					t.Callback()
				}
			}
		}
	}
}

func New(callback func(), interval time.Duration) {
	queue <- &T{
		Interval: interval,
		Callback: callback,
	}
}

func init() {
	queue = make(chan *T)
}
