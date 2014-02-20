package timers

import "time"
import "sync"
import "github.com/tobz/phosphorus/interfaces"

type Timer struct {
	tickInterval  time.Duration
	tickerStop    chan struct{}
	timerSinks    []interfaces.TimerSink
	timerSinkLock *sync.Mutex
}

func NewTimer(d time.Duration) *Timer {
	tickerStop := make(chan struct{}, 1)
	timerSinks := make([]interfaces.TimerSink, 0)
	timerSinkLock := &sync.Mutex{}

	return &Timer{d, tickerStop, timerSinks, timerSinkLock}
}

func (t *Timer) Start() {
	go func() {
		for {
			select {
			case <-t.tickerStop:
				break
			default:
				// Sleep for our tick interval and then notify our sinks.
				time.Sleep(t.tickInterval)

				t.timerSinkLock.Lock()
				for _, sink := range t.timerSinks {
					sink.Tick()
				}
				t.timerSinkLock.Unlock()
			}
		}
	}()
}

func (t *Timer) Stop() {
	t.tickerStop <- struct{}{}
}

func (t *Timer) AddSink(sink interfaces.TimerSink) {
	t.timerSinkLock.Lock()
	defer t.timerSinkLock.Unlock()

	t.timerSinks = append(t.timerSinks, sink)
}
