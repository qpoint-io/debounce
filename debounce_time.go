// Copyright © 2024 Jon Friesen <jon@qpoint.io>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package debounce

import (
	"sync"
	"time"
)

var now = time.Now

// NewDebounceByDuration returns a debounced function that takes another function as its argument.
// This function will be called at the given interval, but no more than the max duration
// from the first call.
func NewDebounceByDuration(interval, maxDuration time.Duration) func(f func()) {
	d := &durationDebouncer{
		interval:    interval,
		maxDuration: maxDuration,
	}

	return func(f func()) {
		d.add(f)
	}
}

type durationDebouncer struct {
	mu          sync.Mutex
	interval    time.Duration
	maxDuration time.Duration
	timer       *time.Timer
	firstCall   bool
	startTime   time.Time
}

func (d *durationDebouncer) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := now()
	if !d.firstCall {
		d.firstCall = true
		d.startTime = now
	}

	if d.timer != nil {
		d.timer.Stop()
	}

	remainingDuration := d.maxDuration - time.Since(d.startTime)
	if remainingDuration <= 0 {
		d.reset()
		f()
		return
	}

	d.timer = time.AfterFunc(d.interval, func() {
		d.mu.Lock()
		defer d.mu.Unlock()

		f()
		d.reset()
	})
}

func (d *durationDebouncer) reset() {
	d.firstCall = false
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
	d.startTime = time.Time{}
}
