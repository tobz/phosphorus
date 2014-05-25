package timers

import "fmt"
import "time"
import "sync"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/statistics"
import "github.com/rcrowley/go-metrics"

type Timer struct {
	tickInterval      time.Duration
	tickerStop        chan struct{}
	timerSinks        []interfaces.TimerSink
	timerSinkLock     *sync.Mutex
	timerOverUnder    metrics.Gauge
	timerTickDuration metrics.Gauge
	timerSinkCount    metrics.Counter
}

func NewTimer(d time.Duration) *Timer {
	tickerStop := make(chan struct{}, 1)
	timerSinks := make([]interfaces.TimerSink, 0)
	timerSinkLock := &sync.Mutex{}

	return &Timer{tickInterval: d, tickerStop: tickerStop, timerSinks: timerSinks, timerSinkLock: timerSinkLock}
}

func NewTrackedTimer(timerName string, d time.Duration) *Timer {
	timer := NewTimer(d)
	timer.timerOverUnder = metrics.GetOrRegisterGauge(fmt.Sprintf("timers.%s.overUnder", timerName), statistics.Registry)
	timer.timerTickDuration = metrics.GetOrRegisterGauge(fmt.Sprintf("timers.%s.tickDuration", timerName), statistics.Registry)
	timer.timerSinkCount = metrics.GetOrRegisterCounter(fmt.Sprintf("timers.%s.sinkCount", timerName), statistics.Registry)

	return timer
}

func (t *Timer) Start() {
	go func() {
		for {
			select {
			case <-t.tickerStop:
				break
			default:
				sleepStart := time.Now()

				// Sleep for our tick interval and then notify our sinks.
				time.Sleep(t.tickInterval)

				if t.timerOverUnder != nil {
					overUnder := (time.Now().Sub(sleepStart).Nanoseconds() - t.tickInterval.Nanoseconds())
					t.timerOverUnder.Update(overUnder)
				}

				tickStart := time.Now()

				t.timerSinkLock.Lock()
				for _, sink := range t.timerSinks {
					sink.Tick()
				}
				t.timerSinkLock.Unlock()

				if t.timerTickDuration != nil {
					tickDuration := time.Now().Sub(tickStart).Nanoseconds()
					t.timerTickDuration.Update(tickDuration)
				}
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

	if t.timerSinkCount != nil {
		t.timerSinkCount.Inc(1)
	}
}
