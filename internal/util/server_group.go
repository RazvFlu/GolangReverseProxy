package util

import "errors"

type ServerGroup struct {
	hasServers bool
	errCh      chan error
}

var ErrNoServers = errors.New("no servers to wait for")
var ErrServerExited = errors.New("server exited without error")

// Call a server function in a goroutine and wait for it to exit.
func NewServerGroup() *ServerGroup {
	return &ServerGroup{
		errCh: make(chan error, 1),
	}
}

func (sg *ServerGroup) Go(f func() error) {
	sg.hasServers = true

	go func() {
		err := f()
		if err == nil {
			err = ErrServerExited
		}

		select {
		case sg.errCh <- err:
		default:
		}
	}()
}

// Wait for all servers to exit.
func (sg *ServerGroup) Wait() error {
	if !sg.hasServers {
		return ErrNoServers
	}

	return <-sg.errCh
}
