package debounce

import (
	"sync"
	"testing"
	"time"
)

var (
	mockNowFunc  func() time.Time
	mockNowMutex sync.Mutex
)

func mockNow() time.Time {
	mockNowMutex.Lock()
	defer mockNowMutex.Unlock()
	return mockNowFunc()
}

func setMockNow(t time.Time) {
	mockNowMutex.Lock()
	defer mockNowMutex.Unlock()
	mockNowFunc = func() time.Time {
		return t
	}
}

func init() {
	now = mockNow
}

func TestTimeDebounce(t *testing.T) {
	tests := []struct {
		name        string
		interval    time.Duration
		maxDuration time.Duration
		timeSteps   []time.Duration
	}{
		{
			name:        "Single call within interval",
			interval:    100 * time.Millisecond,
			maxDuration: 500 * time.Millisecond,
			timeSteps:   []time.Duration{100 * time.Millisecond, 50 * time.Millisecond},
		},
		{
			name:        "Multiple calls within interval",
			interval:    100 * time.Millisecond,
			maxDuration: 500 * time.Millisecond,
			timeSteps:   []time.Duration{50 * time.Millisecond, 50 * time.Millisecond, 200 * time.Millisecond},
		},
		{
			name:        "Single call at max duration",
			interval:    100 * time.Millisecond,
			maxDuration: 200 * time.Millisecond,
			timeSteps:   []time.Duration{100 * time.Millisecond, 100 * time.Millisecond},
		},
		{
			name:        "Multiple calls at max duration",
			interval:    150 * time.Millisecond,
			maxDuration: 300 * time.Millisecond,
			timeSteps:   []time.Duration{50 * time.Millisecond, 50 * time.Millisecond, 200 * time.Millisecond},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			mu := sync.Mutex{}
			f := func() {
				mu.Lock()
				defer mu.Unlock()
				callCount++
			}

			d := NewDebounceByDuration(tt.interval, tt.maxDuration)
			start := time.Now()
			setMockNow(start)

			for _, step := range tt.timeSteps {
				setMockNow(start.Add(step))
				d(f)
				time.Sleep(step)
			}

			mu.Lock()
			if callCount != 1 {
				t.Errorf("expected 1 calls, got %d", callCount)
			}
			mu.Unlock()
		})
	}
}
