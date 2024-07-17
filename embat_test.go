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

	// Set frequency to 50ms to ensure all jobs are processed before shutdown
	mb := embat.NewMicroBatcher[string, int](
		mbp,
		embat.WithFrequency[string, int](1*time.Second),
		embat.WithBatchSize[string, int](1),
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
	assert.True(
		t,
		ts.Seconds() < float64(time.Second*time.Duration(numJobs)),
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
				t.Error("expected result not received in time")
				return
			}
		}
	}()

	wg.Wait()
}
