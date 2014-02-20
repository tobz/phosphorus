package timers

import "time"
import "testing"
import "github.com/stretchr/testify/assert"

type MockSink struct {
	TickCount int
}

func (ms *MockSink) Tick() {
	ms.TickCount += 1
}

func TestTimerSingleSink(t *testing.T) {
	// Create a mock sink.
	mockSink := &MockSink{}

	// Create a timer with a tick interval of one second.
	timer := NewTimer(time.Second * 1)
	timer.AddSink(mockSink)
	timer.Start()

	// Wait two seconds to let the timer tick.
	time.Sleep(time.Second * 2)

	// Stop the timer.
	timer.Stop()

	// Make sure we ticked.
	assert.True(t, mockSink.TickCount > 0, "timer should have ticked at least once")
}

func TestTimerMultiSink(t *testing.T) {
	// Create a mock sink.
	mockSink := &MockSink{}

	// Create a timer with a tick interval of one second.
	// Add the sink to it three times to get 3x the ticks.
	timer := NewTimer(time.Second * 1)
	timer.AddSink(mockSink)
	timer.AddSink(mockSink)
	timer.AddSink(mockSink)
	timer.Start()

	// Wait two seconds to let the timer tick.
	time.Sleep(time.Second * 2)

	// Stop the timer.
	timer.Stop()

	// Make sure we ticked.
	assert.True(t, mockSink.TickCount > 2, "timer should have ticked at least three times")
}
