package interfaces

// Defines an object where it can receives notifications of a world tick i.e.
// a timer that a given time period has elapsed.  It should use this to maintain
// internal state about when it has last run time-bound updates, or updated its
// dependants that require a timer source.
type TimerSink interface {
	Tick()
}
