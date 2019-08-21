package updatingIdentifier

import (
	"sync"
	"time"
)

type Data struct {
	rawValue         uint64
	value, prevValue string
	mutex            sync.RWMutex
	halfTimeout      time.Duration
}

const (
	minTimeout = time.Millisecond * 8
)

func New(timeout time.Duration) *Data {
	if timeout < minTimeout {
		timeout = minTimeout
	}

	i := &Data{
		halfTimeout: timeout / 2,
	}

	i.generateValues()
	go i.watch()

	return i
}
