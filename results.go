package embat

import (
	"sync"
)

// results holds a map of the results of processed jobs.
type results[R any] struct {
	mu sync.Mutex
	m  map[jobID]chan Result[R]
}

// add safely adds a new job result channel to the results map
func (r *results[R]) add(jobID jobID, ch chan Result[R]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[jobID] = ch
}

// sendResults sends the results of processed jobs to the respective result channels.
func (r *results[R]) sendResults(jobResults []Result[R]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, result := range jobResults {
		if ch, ok := r.m[result.JobID]; ok {
			ch <- result
			close(ch)
			delete(r.m, result.JobID)
		}
	}
}
