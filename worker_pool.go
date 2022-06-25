package bits

import (
	"context"
	"sync"
)

type workerPoolConfig struct {
	numOfWorkers int
	buffer       int
}

type WorkerPoolConfigurator func(cfg *workerPoolConfig)

func WorkerPoolNum(num int) WorkerPoolConfigurator {
	return func(cfg *workerPoolConfig) {
		cfg.numOfWorkers = num
	}
}

func WorkerPoolBuffer(buffer int) WorkerPoolConfigurator {
	return func(cfg *workerPoolConfig) {
		cfg.buffer = buffer
	}
}

func WorkerPool[T any, U any](
	ctx context.Context,
	worker func(context.Context, T) (U, error),
	configurators ...WorkerPoolConfigurator,
) (
	chan T,
	chan U,
	chan error,
) {
	cfg := workerPoolConfig{
		numOfWorkers: 4,
		buffer:       999,
	}
	for _, c := range configurators {
		c(&cfg)
	}

	inputChan := make(chan T, cfg.buffer)
	outputChan := make(chan U, cfg.buffer)
	errChan := make(chan error, cfg.buffer)

	wg := sync.WaitGroup{}
	wg.Add(cfg.numOfWorkers)

	go func() {
		wg.Wait()

		close(outputChan)
		close(errChan)
	}()

	for i := 0; i < cfg.numOfWorkers; i++ {
		go func() {
			defer wg.Done()

			var output U
			var err error

			for input := range inputChan {
				output, err = worker(ctx, input)
				if err != nil {
					errChan <- err
				} else {
					outputChan <- output
				}
			}
		}()
	}

	return inputChan, outputChan, errChan
}
