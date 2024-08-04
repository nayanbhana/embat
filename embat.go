//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE -package=mock
package embat

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// BatchProcessor processes a batch of jobs, this interface should be implemented by the consumer.
type BatchProcessor[J any, R any] interface {
	// Process processes a batch of jobs and returns their results.
	Process(batch []Job[J]) []Result[R]
}

// NewMicroBatcher creates a new MicroBatcher with given options.
func NewMicroBatcher[J any, R any](processor BatchProcessor[J, R], opts ...Option[J, R]) *MicroBatcher[J, R] {
	mb := &MicroBatcher[J, R]{
		processor: processor,
		batchSize: 100,
		frequency: 5 * time.Second,
		jobs: jobsC[J]{
			c: make(chan Job[J], 1),
		},
		logger: noOpLogger{},
		results: results[R]{
			m: make(map[JobID]chan Result[R]),
		},
	}

	for _, opt := range opts {
		opt(mb)
	}

	mb.wg.Add(1)
	go mb.start()
	return mb
}

// NewJob creates a new Job with a generated ID.
func NewJob[T any](data T) Job[T] {
	return Job[T]{
		ID:   newJobID(),
		Data: data,
	}
}

// Job represents a unit of work to be processed.
// The generic type J is the specific data for the job supplied by the consumer.
type Job[J any] struct {
	// ID is a unique identifier for the job.
	ID JobID
	// Data holds the specific data for the job, of the generic type J.
	Data J
}

// NewResult creates a new Result with the given JobID and outcome.
func NewResult[R any](jobID JobID, result R, err error) Result[R] {
	return Result[R]{
		JobID:  jobID,
		Result: result,
		Err:    err,
	}
}

// Result represents the result of a processed Job.
// The generic type R is the outcome of processing the job supplied by the consumer.
type Result[R any] struct {
	// JobID is the identifier of the original job that was processed.
	JobID JobID
	// Result holds the outcome of processing the job, of the generic type R.
	Result R
	// Err holds any error that occurred during submission or processing of the job.
	Err error
}

// MicroBatcher handles batching and processing of jobs.
type MicroBatcher[J any, R any] struct {
	// batchSize is the maximum number of jobs in each batch.
	batchSize int
	// frequency is the duration between batch processing attempts.
	frequency time.Duration
	// jobs is the current list of pending jobs to be processed.
	jobs jobsC[J]
	// logger is the logger for the MicroBatcher.
	// default is no logging, if you want logging you can provide your own logger.
	logger Logger
	// processor is the BatchProcessor supplied by the consumer that processes batches of jobs.
	processor BatchProcessor[J, R]
	// results maps each job ID to its result channel.
	results results[R]
	// shutdownOnce ensures that shutdown is called only once.
	shutdownOnce sync.Once
	// shutdownCalled flag to indicate if shutdown has been called.
	shutdownCalled atomic.Bool
	// wg is a wait group to ensure all jobs are processed before shutdown.
	wg sync.WaitGroup
}

// Submit adds a job to the MicroBatcher and returns a channel to receive the result
func (mb *MicroBatcher[J, R]) Submit(job Job[J]) <-chan Result[R] {
	if mb.isShutdown() {
		mb.logger.Debug("shutdown has been initiated, submit failed for job with id: %s", job.ID)
		return mb.shutdownResult(job)
	}
	if job.ID == "" {
		job.ID = newJobID()
	}
	resultCh := make(chan Result[R], 1)
	mb.jobs.add(job)
	mb.results.add(job.ID, resultCh)
	mb.logger.Debug("successfully submitted job with id: %s", job.ID)
	return resultCh
}

// Logger returns the logger for the MicroBatcher.
func (mb *MicroBatcher[J, R]) Logger() Logger {
	return mb.logger
}

// Frequency returns the processing frequency of the MicroBatcher.
func (mb *MicroBatcher[J, R]) Frequency() time.Duration {
	return mb.frequency
}

// BatchSize returns the batch size of the MicroBatcher.
func (mb *MicroBatcher[J, R]) BatchSize() int {
	return mb.batchSize
}

// Shutdown stops the batcher after processing all submitted jobs.
func (mb *MicroBatcher[J, R]) Shutdown() {
	mb.shutdownOnce.Do(func() {
		mb.logger.Debug("shutdown initiated")
		mb.shutdownCalled.Swap(true)
	})
	mb.jobs.close()
	go mb.wg.Wait()
}

// start starts the MicroBatcher and processes jobs in batches.
func (mb *MicroBatcher[J, R]) start() {
	defer mb.wg.Done()
	ticker := time.NewTicker(mb.frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mb.processBatch()
			if mb.isComplete() {
				mb.logger.Debug("all jobs processed, shutting down")
				return
			}
		}
	}
}

// processBatch processes the next batch of jobs and sends the results.
func (mb *MicroBatcher[J, R]) processBatch() {
	batch := mb.jobs.next(mb.batchSize)
	if len(batch) == 0 {
		return
	}
	jobResults := mb.processor.Process(batch)
	mb.results.sendResults(jobResults)
}

// shutdownResult returns a result channel with an error for a job submitted after shutdown.
func (mb *MicroBatcher[J, R]) shutdownResult(job Job[J]) <-chan Result[R] {
	ch := make(chan Result[R], 1)
	ch <- Result[R]{
		JobID: job.ID,
		Err:   errors.New("job submitted after shutdown"),
	}
	close(ch)
	return ch
}

// isShutdown returns true if the MicroBatcher has been shutdown.
func (mb *MicroBatcher[J, R]) isShutdown() bool {
	return mb.shutdownCalled.Load()
}

// isComplete returns true if there are no more jobs to process and shutdown has been called.
func (mb *MicroBatcher[J, R]) isComplete() bool {
	return mb.jobs.length() == 0 && mb.shutdownCalled.Load()
}

type JobID string

func newJobID() JobID {
	return JobID(uuid.New().String())
}
