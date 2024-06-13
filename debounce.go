// Copyright © 2019 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package debounce provides a debouncer func. The most typical use case would be
// the user typing a text into a form; the UI needs an update, but let's wait for
// a break.
package debounce

import (
	"sync"
	"time"
)

// New returns a debounced function that takes another function as its argument.
// This function will be called when the debounced function stops being called
// for the given duration, provided the maximum count hasn't been exceeded.
// Once the maximum count is exceeded, the function is executed one last time
// and the debouncer is reset.
// The debounced function can be invoked with different functions, if needed,
// the last one will win.
func New(after time.Duration, countLimit uint64) func(f func()) {
	d := &debouncer{
		after:      after,
		countLimit: countLimit,
	}

	return func(f func()) {
		d.add(f)
	}
}

type debouncer struct {
	mu         sync.Mutex
	after      time.Duration
	timer      *time.Timer
	count      uint64
	countLimit uint64
}

func (d *debouncer) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Increment the count
	d.count++

	// If count exceeds maxCount, execute the function and reset
	if d.count > d.countLimit {
		if d.timer != nil {
			d.timer.Stop()
			d.timer = nil
		}

		f()

		// Reset the count for the next iteration
		d.count = 0
		return
	}

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.after, func() {
		d.mu.Lock()
		defer d.mu.Unlock()

		f()

		// Reset the count after the function is executed
		d.count = 0
	})
}
