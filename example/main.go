package main

import (
	"fmt"

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
		results[i] = embat.NewResult(job.ID, MyResultType{ProcessedData: fmt.Sprintf("Processed: %s", job.Data.Data)}, nil)
	}
	return results
}

func main() {
	processor := &MyProcessor{}
	batcher := embat.NewMicroBatcher[MyJobType, MyResultType](processor)

	job1 := embat.Job[MyJobType]{Data: MyJobType{Data: "job1"}}

	job2 := embat.NewJob(MyJobType{Data: "job2"})
	job3 := embat.NewJob(MyJobType{Data: "job3"})
	job4 := embat.NewJob(MyJobType{Data: "job4"})
	job5 := embat.NewJob(MyJobType{Data: "job5"})

	resultCh1 := batcher.Submit(job1)
	resultCh2 := batcher.Submit(job2)

	resultCh3 := batcher.Submit(job3)
	resultCh4 := batcher.Submit(job4)
	batcher.Shutdown()
	batcher.Shutdown()
	resultCh5 := batcher.Submit(job5)

	fmt.Printf("%+v\n", <-resultCh1)
	fmt.Printf("%+v\n", <-resultCh2)
	fmt.Printf("%+v\n", <-resultCh3)
	fmt.Printf("%+v\n", <-resultCh4)
	fmt.Printf("%+v\n", <-resultCh5)
	//time.Sleep(10 * time.Second)

}
