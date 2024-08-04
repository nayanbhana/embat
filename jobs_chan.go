package embat

// jobsC holds a slice of jobs to be processed.
type jobsC[J any] struct {
	c chan Job[J]
}

// add safely adds a job to the jobs slice.
func (j *jobsC[J]) add(job Job[J]) {
	j.c <- job
}

// next returns the next batch of jobs to be processed and removes them from the jobs slice.
func (j *jobsC[J]) next(defaultBatchSize int) []Job[J] {
	jLength := len(j.c)
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
	for i := range batch {
		batch[i] = <-j.c
	}

	return batch
}

// length returns the number of jobs in the jobs slice.
func (j *jobsC[J]) length() int {
	return len(j.c)
}

func (j *jobsC[J]) close() {
	close(j.c)
}
