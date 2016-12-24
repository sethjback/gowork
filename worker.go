package gowork

// Worker do work when called by the queue and return their result when finished
type Worker interface {
	// Do the work. Resturn will be passed over the results channel
	Do() interface{}
}
