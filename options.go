package embat

import (
	"time"
)

// Option is a type for configuring the MicroBatcher.
type Option[J any, R any] func(*MicroBatcher[J, R])

// WithFrequency sets the processing frequency.
func WithFrequency[J any, R any](frequency time.Duration) Option[J, R] {
	return func(mb *MicroBatcher[J, R]) {
		mb.frequency = frequency
	}
}

// WithBatchSize sets the batch size.
func WithBatchSize[J any, R any](size int) Option[J, R] {
	return func(mb *MicroBatcher[J, R]) {
		mb.batchSize = size
	}
}

// WithLogger sets the logger for the MicroBatcher.
func WithLogger[J any, R any](logger Logger) Option[J, R] {
	return func(mb *MicroBatcher[J, R]) {
		mb.logger = logger
	}
}
