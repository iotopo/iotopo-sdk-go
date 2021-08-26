package distrlock

import (
	"github.com/bsm/redislock"
	"gogs.iotopo.com/iotopo/iotopo-sdk-go/cache"
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = redislock.ErrNotObtained

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = redislock.ErrLockNotHeld
)
func NewLocker() *redislock.Client {
	return redislock.New(cache.GetClient())
}

