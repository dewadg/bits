package bits

import (
	"context"
	"testing"
)

func TestWorkerPool_Int(t *testing.T) {
	ctx := context.Background()

	inputChan, outputChan, errChan := WorkerPool[int, int](
		ctx,
		func(ctx context.Context, i int) (int, error) {
			return i * i, nil
		},
		WorkerPoolNum(4),
		WorkerPoolBuffer(5),
	)

	go func() {
		for i := 1; i <= 1_000_000; i++ {
			inputChan <- 1
		}
		close(inputChan)
	}()

	var actual int
	for i := range outputChan {
		actual += i
	}
	for range errChan {
	}

	expected := 1_000_000
	if actual != expected {
		t.Errorf("expected result %d, got %d", expected, actual)
	}
}
