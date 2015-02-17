package lob

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

// MetricsBundle is a bundle of a timer and two meters,
// for success and failure. It is useful for, e.g.,
// API endpoints.
type MetricsBundle struct {
	Timer   metrics.Timer
	Success metrics.Meter
	Error   metrics.Meter
}

// NewMetricsBundle creates a new metrics bundle and
// registers them with the given name as a prefix.
func NewMetricsBundle(name string) *MetricsBundle {
	m := &MetricsBundle{
		Timer:   metrics.NewTimer(),
		Success: metrics.NewMeter(),
		Error:   metrics.NewMeter(),
	}
	err := metrics.Register(name+".timer", m.Timer)
	if err != nil {
		panic(err)
	}
	metrics.Register(name+".success", m.Success)
	if err != nil {
		panic(err)
	}
	metrics.Register(name+".error", m.Error)
	if err != nil {
		panic(err)
	}
	return m
}

// Call calls a function returning an error, and passes
// it back, while recording timing and success information.
func (m *MetricsBundle) Call(f func() error) error {
	ts := time.Now()
	defer m.Timer.UpdateSince(ts)

	err := f()
	if err != nil {
		m.Error.Mark(1)
	} else {
		m.Success.Mark(1)
	}
	return err
}
