# gowork

gowork is a drop in work queue for running jobs on a configurable number of go routines.

# Usage

When initializing the queue, you specify the depth of both the incoming work and the outgoing results. This is one way to fine tune the execution rate.

```go
// queue with an incoming work buffer of 100, and an outgoing result buffer of 10.
q := NewQueue(100, 10)
```

After the queue has been created, you start it specifying the number of go routines you want to execute work simultaneously

```go
// queue started with 4 go routines
q.Start(4)
```

You add work to the queue by passing something that implements the Worker interface. It will return after queuing the work. Once the work buffer limit has been reached, this will block until there is room to add the work to the queue.

It important that you call Finish() once you are done adding work to the queue. This will signal the queue that no more work will be added, and it should drain the queue then wrap up. Once you call Finish(), you must NOT add more work to the queue.

```go
q.AddWork(&worker{})

q.Finish()
```

As work finishes, the results will be sent over the result channel. You can range over this channel to process the results. The channel will be closed once all the work has been completed.

```go

for r := range q.Results() {
  // process results
}
```

If for some reason you need to stop processing the queue, calling Abort() will allow the existing work to finish, but not process anymore. Once the current work has finished, the results channel will be closed as normal

```go
// Will NOT process any more queued work
q.Abort()
```
