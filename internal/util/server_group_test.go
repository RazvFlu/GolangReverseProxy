package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerGroupEmpty(t *testing.T) {
	sg := NewServerGroup()
	assert.ErrorIs(t, sg.Wait(), ErrNoServers)
}

func TestServerGroupSingle(t *testing.T) {
	sg := NewServerGroup()

	sg.Go(func() error { return nil })

	assert.Error(t, sg.Wait(), ErrServerExited)
}

func TestServerGroupMultiple(t *testing.T) {
	sg := NewServerGroup()

	for i := 0; i < 10; i++ {
		sg.Go(func() error { return nil })
	}

	assert.Error(t, sg.Wait(), ErrServerExited)
}

func TestServerGroupError(t *testing.T) {
	sg := NewServerGroup()

	sg.Go(func() error { return ErrServerExited })

	assert.ErrorIs(t, sg.Wait(), ErrServerExited)
}
