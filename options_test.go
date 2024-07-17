package embat_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nayanbhana/embat"
	"github.com/nayanbhana/embat/mock"
)

// TestWithFrequency tests that the frequency option can be set.
func TestWithFrequency(t *testing.T) {
	f := time.Second * 100
	mb := embat.NewMicroBatcher[int, int](nil, embat.WithFrequency[int, int](f))
	assert.Equal(t, f, mb.Frequency())
}

// TestWithSize tests that the size option can be set.
func TestWithSize(t *testing.T) {
	bs := 1000
	mb := embat.NewMicroBatcher[int, int](nil, embat.WithBatchSize[int, int](bs))
	assert.Equal(t, bs, mb.BatchSize())
}

// TestWithLogger tests that the logger option can be set.
func TestWithLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l := mock.NewMockLogger(ctrl)
	mb := embat.NewMicroBatcher[int, int](nil, embat.WithLogger[int, int](l))
	assert.Equal(t, l, mb.Logger())
}
