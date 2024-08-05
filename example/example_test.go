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
	Result string
}

type MyProcessor struct{}

func (p *MyProcessor) Process(batch []embat.Job[MyJobType]) []embat.Result[MyResultType] {
	results := make([]embat.Result[MyResultType], len(batch))
	for i, job := range batch {
		results[i] = embat.NewResult(
			job.ID,
			MyResultType{Result: fmt.Sprintf("successfully processed %s", job.Data.Data)},
			nil,
		)
	}
	return results
}

func TestExample(t *testing.T) {
	numberOfJobs := 1000
	batchSize := 100
	frequency := time.Second * 1

	// Define consumer processor.
	processor := &MyProcessor{}

	// Create a new MicroBatcher.
	batcher := embat.NewMicroBatcher[MyJobType, MyResultType](
		processor,
		embat.WithFrequency[MyJobType, MyResultType](frequency),
		embat.WithBatchSize[MyJobType, MyResultType](batchSize),
		embat.WithLogger[MyJobType, MyResultType](&logger{t, true}),
	)

	wg := &sync.WaitGroup{}
	results := make(chan embat.Result[MyResultType], 1)

	// Create and submit jobs.
	for i := 1; i <= numberOfJobs; i++ {
		j := createJob(fmt.Sprintf("job number: %v", i))
		wg.Add(1)
		go submitJob(batcher, wg, results, j)
	}

	// Wait for all jobs to complete.
	go wait(wg, results)

	// Print results as jobs are done
	printResults(results)

	// Shutdown the batcher.
	batcher.Shutdown()

	// Submit a job after shutdown.
	afterShutdown := batcher.Submit(createJob("job after shutdown"))
	fmt.Printf("%+v\n", <-afterShutdown)
}

func submitJob(batcher *embat.MicroBatcher[MyJobType, MyResultType], wg *sync.WaitGroup, results chan embat.Result[MyResultType], job embat.Job[MyJobType]) {
	results <- <-batcher.Submit(job)
	wg.Done()
}

func wait(wg *sync.WaitGroup, results chan embat.Result[MyResultType]) {
	defer close(results)
	wg.Wait()
}

func printResults(results chan embat.Result[MyResultType]) {
	for r := range results {
		fmt.Printf("%+v\n", r)
	}
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
