package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nayanbhana/embat"
)

type MyJobType struct {
	Data string
}

type MyResultType struct {
	ProcessedData string
}

type MyProcessor struct{}

func (p *MyProcessor) Process(batch []embat.Job[MyJobType]) []embat.Result[MyResultType] {
	results := make([]embat.Result[MyResultType], len(batch))
	for i, job := range batch {
		results[i] = embat.NewResult(
			job.ID,
			MyResultType{ProcessedData: fmt.Sprintf("Processed: %s", job.Data.Data)},
			nil,
		)
	}
	return results
}

func TestExample(t *testing.T) {
	// Define consumer processor.
	processor := &MyProcessor{}

	// Create a new MicroBatcher.
	batcher := embat.NewMicroBatcher[MyJobType, MyResultType](
		processor,
		embat.WithFrequency[MyJobType, MyResultType](time.Second*1),
		embat.WithBatchSize[MyJobType, MyResultType](2),
		embat.WithLogger[MyJobType, MyResultType](&logger{t, true}),
	)

	numberOfJobs := 6

	wg := sync.WaitGroup{}
	results := make(chan embat.Result[MyResultType], numberOfJobs)

	// Create and submit jobs.
	for i := 1; i <= numberOfJobs; i++ {
		wg.Add(1)
		go func(i int) {
			results <- <-batcher.Submit(createJob(fmt.Sprintf("job number: %v", i)))
		}(i)
	}

	// Wait for all jobs to be processed.
	go func() {
		defer close(results)
		wg.Wait()
	}()

	// Print results.
	for r := range results {
		fmt.Printf("%+v\n", r)
		wg.Done()
	}

	// Shutdown the batcher.
	batcher.Shutdown()

	// Submit a job after shutdown.
	afterShutdown := batcher.Submit(createJob("job after shutdown"))
	fmt.Printf("%+v\n", <-afterShutdown)
}

func createJob(data string) embat.Job[MyJobType] {
	return embat.NewJob(MyJobType{Data: data})
}

type logger struct {
	t       testing.TB
	enabled bool
}

func (f *logger) Debug(format string, args ...any) {
	if f.enabled {
		f.t.Logf(format, args...)
	}
}
