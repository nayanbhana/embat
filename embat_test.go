package embat_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nayanbhana/embat"
	"github.com/nayanbhana/embat/mock"
)

// TestMicroBatcher_Submit tests that the MicroBatcher processes jobs correctly.
func TestMicroBatcher_Submit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mbp := mock.NewMockBatchProcessor[string, int](ctrl)
	mbp.EXPECT().
		Process(gomock.Any()).
		DoAndReturn(func(jobs []embat.Job[string]) []embat.Result[int] {
			var results []embat.Result[int]
			for _, job := range jobs {
				r := embat.NewResult(job.ID, 42, nil)
				results = append(results, r)
			}
			return results
		}).Times(1)

	mb := embat.NewMicroBatcher[string, int](
		mbp,
		embat.WithFrequency[string, int](50*time.Millisecond),
		embat.WithBatchSize[string, int](2),
	)

	job := embat.NewJob("test-job")
	resultCh := mb.Submit(job)

	select {
	case result := <-resultCh:
		assert.NoError(t, result.Err)
		assert.Equal(t, 42, result.Result)
	case <-time.After(100 * time.Millisecond):
		t.Error("expected result not received in time")
	}

	mb.Shutdown()
}

// TestMicroBatcher_Shutdown tests that no new jobs are accepted after shutdown.
func TestMicroBatcher_Shutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mbp := mock.NewMockBatchProcessor[string, int](ctrl)
	mb := embat.NewMicroBatcher[string, int](
		mbp,
		embat.WithFrequency[string, int](50*time.Millisecond),
		embat.WithBatchSize[string, int](2),
	)

	mb.Shutdown()
	newJob := embat.NewJob("test-job-after-shutdown")
	newResultCh := mb.Submit(newJob)

	select {
	case result := <-newResultCh:
		assert.Error(t, result.Err)
	case <-time.After(100 * time.Millisecond):
		t.Error("expected result not received in time")
	}
}

// TestMicroBatcher_Shutdown_completes_all_jobs tests thats all jobs are completed after shutdown.
func TestMicroBatcher_Shutdown_completes_all_jobs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mbp := mock.NewMockBatchProcessor[string, int](ctrl)
	mbp.EXPECT().
		Process(gomock.Any()).
		DoAndReturn(func(jobs []embat.Job[string]) []embat.Result[int] {
			var results []embat.Result[int]
			for _, job := range jobs {
				r := embat.NewResult(job.ID, 42, nil)
				results = append(results, r)
			}
			return results
		}).Times(3)

	frequency := time.Second
	batchSize := 1
	mb := embat.NewMicroBatcher[string, int](
		mbp,
		embat.WithFrequency[string, int](frequency),
		embat.WithBatchSize[string, int](batchSize),
	)

	// Number of jobs to submit
	numJobs := 3
	var wg sync.WaitGroup
	wg.Add(numJobs)

	results := make([]<-chan embat.Result[int], numJobs)
	n := time.Now()
	for i := 0; i < numJobs; i++ {
		job := embat.NewJob(fmt.Sprintf("test-job-%v", i))
		resultCh := mb.Submit(job)
		results[i] = resultCh
	}
	mb.Shutdown()
	ts := time.Since(n)

	// Assert that the shutdown happened before all jobs were processed
	// The time to process 3 jobs at 1 second each would be 3 seconds
	assert.True(
		t,
		ts < 3*time.Second,
		"shutdown before all jobs are processed",
	)

	// Collect results
	go func() {
		for _, resultCh := range results {
			select {
			case result := <-resultCh:
				assert.NoError(t, result.Err)
				assert.Equal(t, 42, result.Result)
				wg.Done()
			case <-time.After(15 * time.Second):
				assert.Fail(t, "expected result not received in time")
			}
		}
	}()

	wg.Wait()
}
