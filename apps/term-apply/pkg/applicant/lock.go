package applicant

import "sync"

type LockVendor struct {
	locks sync.Map
}

func (self *LockVendor) LockForName(name string) *sync.Mutex {
	lock, _ := self.locks.LoadOrStore(name, &sync.Mutex{})
	if lock, ok := lock.(*sync.Mutex); ok {
		return lock
	}
	panic("Failed to store mutex")
}

func NewLockVendor() *LockVendor {
	return &LockVendor{}
}
