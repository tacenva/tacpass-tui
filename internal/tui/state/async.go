package state

type Async struct {
	Loading bool
	Error   error
}

func (s *Async) Start() {
	s.Loading = true
	s.Error = nil
}

func (s *Async) Success() {
	s.Loading = false
	s.Error = nil
}

func (s *Async) Fail(err error) {
	s.Loading = false
	s.Error = err
}

func (s *Async) Reset() {
	s.Loading = false
	s.Error = nil
}
