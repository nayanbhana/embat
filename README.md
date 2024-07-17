![embat Logo](./logo.svg)

Embat is a Go library for handling batching and processing of jobs using a micro-batching approach.

## Features

- **Micro-Batching**: Efficiently processes jobs in batches to optimize resource usage.
- **Customizable**: Allows customization of batch size, processing frequency, and logging.
- **Error Handling**: Handles job processing errors and shutdown gracefully.

## Requirements

- Go v1.18 or higher (for generics support)

## Installation

Install Embat using `go get`:

```sh
go get github.com/nayanbhana/embat
```

## Run example and tests

```sh
make example
make test
```

## Usage

### 1. Define a BatchProcessor

Implement the `BatchProcessor` interface to define how jobs are processed:

```go
type BatchProcessor[J any, R any] interface {
    Process(batch []Job[J]) []Result[R]
}
```

### 2. Create Jobs and Results

Define your job and result types using the generic `Job` and `Result` types:

```go
type Job[J any] struct {
    ID   jobID
    Data J
}

type Result[R any] struct {
    JobID jobID
    Result R
    Err    error
}
```

### 3. Initialize a MicroBatcher

Check out an example in the example folder.

### 4. Options

#### WithFrequency

Default frequency is set to 5 seconds.
Sets the processing interval for the batcher. Example:

```go
embat.WithFrequency[J, R](1 * time.Second)
```

#### WithBatchSize

Default batch size is set to 100.
Sets the batch size for processing. Example:

```go
embat.WithBatchSize[J, R](10)
```

#### WithLogger

Sets a custom logger for the MicroBatcher. 
Logging is disabled by default unless a logger is provided.

```go
type Logger interface {
	// Debug prints a debug message.
	Debug(format string, args ...any)
}
```

Example:

```go
embat.WithLogger[R, J](customLogger)
```

## Contributing

Feel free to contribute by submitting issues and pull requests on GitHub at [github.com/nayanbhana/embat](https://github.com/nayanbhana/embat).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

