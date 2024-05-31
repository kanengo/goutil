package dislock

import "errors"

var ErrBlocking = errors.New("dislock: blocking")

type DisLocker interface {
	Lock() error
	TryLock() (bool, error)
	Unlock() error
}
