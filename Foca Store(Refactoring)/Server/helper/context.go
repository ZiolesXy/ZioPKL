package helper

import "time"

// SimpleContext untuk background operations
type SimpleContext struct {
	done chan struct{}
}

func (s *SimpleContext) Done() <-chan struct{} {
	return s.done
}

func (s *SimpleContext) Err() error {
	return nil
}

func (s *SimpleContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (s *SimpleContext) Value(key interface{}) interface{} {
	return nil
}

func Background() *SimpleContext {
	return &SimpleContext{done: make(chan struct{})}
}