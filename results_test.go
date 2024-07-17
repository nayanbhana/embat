package embat

import (
	"sync"
	"testing"
)

// Test_results_add tests the add method of the results type.
func Test_results_add(t *testing.T) {
	// Initialise results type.
	var r results[int]
	r.m = make(map[jobID]chan Result[int])
	const numJobs = 5

	// Add channels concurrently.
	var wg sync.WaitGroup
	wg.Add(numJobs)
	for i := 0; i < numJobs; i++ {
		go func(id int) {
			defer wg.Done()
			jobID := jobID(rune(id))
			ch := make(chan Result[int], 1)
			r.add(jobID, ch)
		}(i)
	}
	wg.Wait()

	// Check if all channels are added
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.m) != numJobs {
		t.Errorf("Expected %d channels, got %d", numJobs, len(r.m))
	}
}

func Test_results_sendResults(t *testing.T) {
	// Initialise results type.
	var r results[int]
	r.m = make(map[jobID]chan Result[int])
	const numJobs = 3

	// Add channels to results.
	for i := 0; i < numJobs; i++ {
		jobID := jobID(rune(i))
		ch := make(chan Result[int], 1)
		r.add(jobID, ch)
	}

	// Prepare job results.
	var jobResults []Result[int]
	for i := 0; i < numJobs; i++ {
		jobID := jobID(rune(i))
		jobResults = append(jobResults, Result[int]{JobID: jobID, Result: i * 2})
	}

	// Send results to channels.
	r.sendResults(jobResults)

	// Check if all channels are closed and removed from results.
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.m) != 0 {
		t.Errorf("Expected all channels to be removed, found %d channels remaining", len(r.m))
	}
}
