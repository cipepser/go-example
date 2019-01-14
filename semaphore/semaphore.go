package semaphore

import (
	"errors"
	"time"
)

var (
	ErrNoTickets      = errors.New("semaphore: could not acquire semaphore")
	ErrIllegalRelease = errors.New("semaphore: can't release the semaphore without acquiring it first")
)

type Interface interface {
	Acquire() error
	Release() error
}

type Implementation struct {
	sem     chan struct{}
	timeout time.Duration
}

func (s *Implementation) Acquire() error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-time.After(s.timeout):
		return ErrNoTickets
	}
}

func (s *Implementation) Release() error {
	select {
	case <-s.sem:
		return nil
	case <-time.After(s.timeout):
		return ErrIllegalRelease
	}
}

func NewInterface(tickets int, timeout time.Duration) Interface {
	return &Implementation{
		sem:     make(chan struct{}, tickets),
		timeout: timeout,
	}
}
