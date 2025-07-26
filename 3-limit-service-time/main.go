//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"context"
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

var (
	userProcessQuota = make(map[int]int)
	mx               = sync.Mutex{}
	processLimit     = 10
)

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	mx.Lock()
	if _, ok := userProcessQuota[u.ID]; !ok {
		userProcessQuota[u.ID] = processLimit
	}
	mx.Unlock()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		process()
	}()

	for {
		select {
		case <-time.Tick(time.Second):
			if !u.IsPremium {
				mx.Lock()
				userProcessQuota[u.ID]--
				leftQuota := userProcessQuota[u.ID]
				mx.Unlock()

				if leftQuota <= 0 {
					return false
				}
			}

		case <-ctx.Done():
			return true
		}
	}
}

func main() {
	RunMockServer()
}
