package gowork

// Worker does the work
type Worker interface {
	Do() interface{} // Do something. The return is passed over the Results channel
}
