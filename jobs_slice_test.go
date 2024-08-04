package embat

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_jobID tests that jobs are being added to the jobs struct.
func Test_jobs_add(t *testing.T) {
	var j jobsS[int]
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
	j.mu.Lock()
	defer j.mu.Unlock()
	if len(j.s) != numJobs {
		t.Errorf("Expected %d jobs, got %d", numJobs, len(j.s))
	}
}

// Test_jobs_next tests that the next batch of jobs is returned and removed from the jobs struct.
func Test_jobs_next(t *testing.T) {
	type testCase[J any] struct {
		name             string
		defaultBatchSize int
		j                jobsS[J]
		nextBatch        []Job[J]
		remainingJobs    []Job[J]
	}
	tests := []*testCase[int]{
		{
			name:             "empty jobs",
			defaultBatchSize: 3,
			j: func() jobsS[int] {
				return jobsS[int]{}
			}(),
			nextBatch:     []Job[int]{},
			remainingJobs: []Job[int](nil),
		},
		{
			name:             "less jobs than batch size",
			defaultBatchSize: 3,
			j: func() jobsS[int] {
				return jobsS[int]{
					s: []Job[int]{
						{ID: JobID('a'), Data: 1},
					},
				}
			}(),
			nextBatch: []Job[int]{
				{ID: JobID('a'), Data: 1},
			},
			remainingJobs: []Job[int]{},
		},
		{
			name:             "equal jobs and batch size",
			defaultBatchSize: 1,
			j: func() jobsS[int] {
				return jobsS[int]{
					s: []Job[int]{
						{ID: JobID('a'), Data: 1},
					},
				}
			}(),
			nextBatch: []Job[int]{
				{ID: JobID('a'), Data: 1},
			},
			remainingJobs: []Job[int]{},
		},
		{
			name:             "more jobs than batch size",
			defaultBatchSize: 1,
			j: func() jobsS[int] {
				return jobsS[int]{
					s: []Job[int]{
						{ID: JobID('a'), Data: 1},
						{ID: JobID('b'), Data: 2},
					},
				}
			}(),
			nextBatch: []Job[int]{
				{ID: JobID('a'), Data: 1},
			},
			remainingJobs: []Job[int]{
				{ID: JobID('b'), Data: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.nextBatch, tt.j.next(tt.defaultBatchSize), "next(%v)", tt.defaultBatchSize)
			assert.Equalf(t, tt.remainingJobs, tt.j.s, "next(%v)", tt.defaultBatchSize)
		})
	}
}
