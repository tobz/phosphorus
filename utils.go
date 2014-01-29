package phosphorus

import "sync"

type Coordinator struct {
	stop      chan bool
	done      chan bool
	input     chan bool
	lock      *sync.Mutex
	listeners []chan bool
}

func NewCoordinator() *Coordinator {
	coordinator := &Coordinator{
		stop:      make(chan bool, 1),
		done:      make(chan bool, 1),
		input:     make(chan bool, 1),
		lock:      &sync.Mutex{},
		listeners: make([]chan bool, 1),
	}

	go func(c *Coordinator) {
		for {
			select {
			case <-c.stop:
				c.done <- true
				break
			case <-c.input:
				c.lock.Lock()
				for _, cc := range c.listeners {
					cc <- true
				}
				c.lock.Unlock()
			}
		}
	}(coordinator)

	return coordinator
}

func (c *Coordinator) Register() <-chan bool {
	listener := make(chan bool, 1)
	c.lock.Lock()
	c.listeners = append(c.listeners, listener)
	c.lock.Unlock()

	return listener
}

func (c *Coordinator) Ping() {
	c.input <- true
}

func (c *Coordinator) Stop() {
	c.stop <- true
	<-c.done
}
