package main

import (
	"time"

	"github.com/cipepser/go-example/observerPattern"
)

func main() {
	en := observerPattern.NewEventNotifier()

	en.Register(observerPattern.NewEventObserver(1))
	en.Register(observerPattern.NewEventObserver(2))

	stop := time.NewTimer(10 * time.Second).C
	tick := time.NewTicker(time.Second).C

	for {
		select {
		case <-stop:
			return
		case t := <-tick:
			en.Notify(observerPattern.Event{
				Date: t.UnixNano(),
			})
		}
	}
}
