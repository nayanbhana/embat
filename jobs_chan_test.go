package embat

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_jobsC_add tests that jobs are being added to the jobs struct.
func Test_jobsC_add(t *testing.T) {
	j := jobsC[int]{c: make(chan Job[int], 15)}
	var wg sync.WaitGroup
	const numJobs = 10

	wg.Add(numJobs)
	// Add jobs concurrently
	for i := 0; i < numJobs; i++ {
		go func(id int) {
			defer wg.Done()
			j.add(Job[int]{ID: JobID(rune(id)), Data: id})
		}(i)
	}
	wg.Wait()

	// Check length of jobs.s
	if j.length() != numJobs {
		t.Errorf("Expected %d jobs, got %d", numJobs, len(j.c))
	}
}

// Test_jobsC_next tests that the next batch of jobs is returned and removed from the jobs struct.
func Test_jobsC_next(t *testing.T) {
	type testCase[J any] struct {
		name                  string
		defaultBatchSize      int
		j                     jobsC[J]
		nextBatch             []Job[J]
		numberOfRemainingJobs int
	}
	tests := []*testCase[int]{
		{
			name:             "empty jobs",
			defaultBatchSize: 3,
			j: func() jobsC[int] {
				return jobsC[int]{}
			}(),
			nextBatch:             []Job[int]{},
			numberOfRemainingJobs: 0,
		},
		{
			name:             "less jobs than batch size",
			defaultBatchSize: 3,
			j: func() jobsC[int] {
				j := jobsC[int]{
					c: make(chan Job[int], 3),
				}
				j.add(Job[int]{ID: JobID('a'), Data: 1})
				return j
			}(),
			nextBatch: []Job[int]{
				{ID: JobID('a'), Data: 1},
			},
			numberOfRemainingJobs: 0,
		},
		{
			name:             "equal jobs and batch size",
			defaultBatchSize: 1,
			j: func() jobsC[int] {
				j := jobsC[int]{
					c: make(chan Job[int], 1),
				}
				j.add(Job[int]{ID: JobID('a'), Data: 1})
				return j
			}(),
			nextBatch: []Job[int]{
				{ID: JobID('a'), Data: 1},
			},
			numberOfRemainingJobs: 0,
		},
		{
			name:             "more jobs than batch size",
			defaultBatchSize: 1,
			j: func() jobsC[int] {
				j := jobsC[int]{
					c: make(chan Job[int], 2),
				}
				j.add(Job[int]{ID: JobID('a'), Data: 1})
				j.add(Job[int]{ID: JobID('b'), Data: 2})
				return j
			}(),
			nextBatch: []Job[int]{
				{ID: JobID('a'), Data: 1},
			},
			numberOfRemainingJobs: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.nextBatch, tt.j.next(tt.defaultBatchSize), "next(%v)", tt.defaultBatchSize)
			assert.Equalf(t, tt.numberOfRemainingJobs, len(tt.j.c), "next(%v)", tt.defaultBatchSize)
		})
	}
}
