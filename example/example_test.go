package main

import (
	"fmt"
	"testing"

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
	processor := &MyProcessor{}
	batcher := embat.NewMicroBatcher[MyJobType, MyResultType](
		processor,
		embat.WithFrequency[MyJobType, MyResultType](1),
		embat.WithBatchSize[MyJobType, MyResultType](2),
		embat.WithLogger[MyJobType, MyResultType](&logger{t, true}),
	)

	job1 := embat.NewJob(MyJobType{Data: "job1"})
	job2 := embat.NewJob(MyJobType{Data: "job2"})
	job3 := embat.NewJob(MyJobType{Data: "job3"})
	job4 := embat.NewJob(MyJobType{Data: "job4"})
	job5 := embat.NewJob(MyJobType{Data: "job5"})

	resultCh1 := batcher.Submit(job1)
	resultCh2 := batcher.Submit(job2)
	resultCh3 := batcher.Submit(job3)
	resultCh4 := batcher.Submit(job4)
	resultCh5 := batcher.Submit(job5)
	batcher.Shutdown()

	fmt.Printf("%+v\n", <-resultCh1)
	fmt.Printf("%+v\n", <-resultCh2)
	fmt.Printf("%+v\n", <-resultCh3)
	fmt.Printf("%+v\n", <-resultCh4)
	fmt.Printf("%+v\n", <-resultCh5)

}

// logger is a Logger for unit tests
type logger struct {
	t       testing.TB
	enabled bool
}

func (f *logger) Debug(format string, args ...any) {
	if f.enabled {
		f.t.Logf(format, args...)
	}
}
