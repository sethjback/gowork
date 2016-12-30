package gowork

import "sync"

// Queue is the collection point for pending work.
// It manages running incoming work on go routines, stops them
// if aborted, and cleaning up after they are finished
type Queue struct {
	workQueue  chan Worker
	resultChan chan interface{}
	abortChan  chan struct{}
}

// NewQueue returns a new work queue.
// workDepth sets the buffer for incoming work: the queue will accept
// up to this number of pending jobs before AddWork will block.
// resultDepth sets the output buffer. Once this fills work will stop
// processing until results are consumed via the Results channel
func NewQueue(workDepth, resultDepth int) *Queue {
	return &Queue{
		workQueue:  make(chan Worker, workDepth),
		resultChan: make(chan interface{}, resultDepth),
		abortChan:  make(chan struct{}),
	}
}

// Results returns a channel over which all worker results are passed.
// Gowork does not distinguish between success and error results, it
// is up to the producer and consumer to decide how to communicate.
func (q *Queue) Results() <-chan interface{} {
	return q.resultChan
}

// AddWork to the queue. Once the input buffer has been filled
// this will block until space opens up in the buffer again.
// Calling this after calling Finish will result in a panic
func (q *Queue) AddWork(w Worker) {
	q.workQueue <- w
}

// Abort will signal all workers to exit after finishing their current task
// This aborts processing the queue: any unfinished work in the buffer will be ignored.
func (q *Queue) Abort() {
	close(q.abortChan)
}

// Finish tells the queue that no more work will be added and the queue should
// complete the work in its buffer and then stop. You must call this
// when you are through with the queue otherwise the worker go routines will continue
// to listen for work. It is safe to call before all work is complete: workers
// will process any backlog in the work queue before exiting
func (q *Queue) Finish() {
	close(q.workQueue)
}

// Start spins up workerCount go routines to process work concurrently.
func (q *Queue) Start(workerCount int) {
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			doWork(q.workQueue, q.abortChan, q.resultChan)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(q.resultChan)
	}()
}

func doWork(q <-chan Worker, abort <-chan struct{}, resultChan chan<- interface{}) {
	for work := range q {
		select {
		case resultChan <- work.Do():
		case <-abort:
			return
		}
	}
}
