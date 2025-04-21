package job

import (
	"kafkador/cmd/service"
	"time"
)

type CleanUpJob struct {
}

func startCleanupLoop(service *service.BikerServiceImpl, interval time.Duration, expiry time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			service.ResponseChans.Range(func(key, value any) bool {
				if wrapper, ok := value.(*service.ResponseWrapper); ok {
					if time.Since(wrapper.CreatedAt) > expiry {
						service.ResponseChans.Delete(key)
						close(wrapper.Chan) // optional, only if no one is listening
					}
				}
				return true
			})
		}
	}()
}
