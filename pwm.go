package pwm

import (
	"sync"
	"time"
)

// An OutputLine represents the gpio pin to toggle
type OutputLine interface {
	SetValue(val int) error
}

// A SoftPWM is a software pulse-width modulation controller which
// controls exactly one OutputLine
type SoftPWM struct {
	dutyCycle uint32
	stop      chan struct{}
	done      chan struct{}
	l         OutputLine
	running   bool
	mu        sync.Mutex
}

// New creates a new SoftPWM. The provided output line is turned off
// initially, and the output duty cycle is set to 0.
func New(l OutputLine) *SoftPWM {
	l.SetValue(0)
	return &SoftPWM{
		l:    l,
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
}

// Set sets the output duty cycle of the OutputLine to a value between
// zero and 1000, inclusive. The minimum pulse width is 100 microseconds.
// Combined with the range of 0-1000, this gives a frequency of ~10Hz.
func (p *SoftPWM) Set(v int) {
	if v < 0 || v > 1000 {
		panic("pwm value out of bounds")
	}
	//atomic.StoreUint32(&p.dutyCycle, uint32(v))
	p.dutyCycle = uint32(v)
	p.mu.Lock()
	if !p.running {
		go p.run()
	}
	p.mu.Unlock()
}

// Off sets the pulse-width to zero and stops the processing loop
func (p *SoftPWM) Off() {
	p.mu.Lock()
	if p.running {
		p.stop <- struct{}{}
		<-p.done
	}
	p.mu.Unlock()
}

func (p *SoftPWM) run() {
	p.running = true
	t := time.NewTicker(100 * time.Microsecond)
	//t := time.NewTicker(1 * time.Millisecond)

	var i uint32

LOOP:
	for {
		select {
		case <-t.C:
			//v := atomic.LoadUint32(&p.dutyCycle)
			if i < p.dutyCycle {
				p.l.SetValue(1)
			} else {
				p.l.SetValue(0)
			}
			i = (i + 1) % 1000
		case <-p.stop:
			p.running = false
			break LOOP
		}
	}
	p.l.SetValue(0)
	p.done <- struct{}{}
}
