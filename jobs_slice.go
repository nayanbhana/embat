package embat

import (
	"sync"
)

// jobsS holds a slice of jobs to be processed.
type jobsS[J any] struct {
	mu sync.Mutex
	s  []Job[J]
}

// add safely adds a job to the jobs slice.
func (j *jobsS[J]) add(job Job[J]) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.s = append(j.s, job)
}

// next returns the next batch of jobs to be processed and removes them from the jobs slice.
func (j *jobsS[J]) next(defaultBatchSize int) []Job[J] {
	j.mu.Lock()
	defer j.mu.Unlock()

	jLength := len(j.s)
	// If there are no jobs, return an empty slice.
	if jLength == 0 {
		return []Job[J]{}
	}

	batchSize := defaultBatchSize
	// If the number of jobs is less than the default batch size, process all jobs.
	if jLength < defaultBatchSize {
		batchSize = jLength
	}
	batch := make([]Job[J], batchSize)
	// Copy the jobs to be processed the batch.
	copy(batch, j.s[:batchSize])
	// Remove the jobs to be processed from the slice.
	j.s = j.s[batchSize:]
	return batch
}

// length returns the number of jobs in the jobs slice.
func (j *jobsS[J]) length() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return len(j.s)
}

// close fulfills the interface but does nothing for this implementation.
func (j *jobsS[J]) close() {
	return
}
