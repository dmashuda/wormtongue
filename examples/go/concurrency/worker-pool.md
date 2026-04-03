# Worker Pool Pattern

Distributes work across a fixed number of goroutines for bounded concurrency.

## When to Use

- Processing a large batch of independent tasks
- Rate-limiting concurrent operations (e.g., API calls, DB writes)
- CPU-bound work that benefits from parallelism up to GOMAXPROCS

## Example

```go
package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker %d processing job %d\n", id, j)
		results <- j * 2
	}
}

func main() {
	const numWorkers = 3
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println("result:", r)
	}
}
```

## Key Points

- Fixed goroutine count prevents unbounded concurrency
- Channel-based work distribution is idiomatic Go
- WaitGroup coordinates clean shutdown
- Close the jobs channel to signal workers there's no more work
- A separate goroutine waits for all workers before closing results
