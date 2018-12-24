// package ratelimit provides the ratelimit implement by token bucket
package ratelimit

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type tokenBucket struct {
	TokenBucket chan uuid.UUID
}

// create a new token bucket with given capacity
func InitTokenBucket(capacity int32) *tokenBucket {
	ch := make(chan uuid.UUID, capacity)
	return &tokenBucket{
		TokenBucket: ch,
	}
}

// preheat a token bucket before it start to fill token
func (t *tokenBucket) Preheat(reserved int32) error {
	if reserved > int32(cap(t.TokenBucket)) {
		return errors.Errorf("reserved:%d shall not bigger than tokenBucket capacity:%d", reserved, cap(t.TokenBucket))
	}
	if len(t.TokenBucket) != 0 {
		return errors.Errorf("preheat shall only used for empty bucket!")
	}
	for i := int32(0); i < reserved; i++ {
		uid, _ := uuid.NewV4()
		select {
		case t.TokenBucket <- uid:
		default:
		}
	}
	return nil
}

// flush the token bucket during the given time duration
func (t *tokenBucket) Flush(flushDuration time.Duration) error {
	ticker := time.NewTicker(flushDuration)
	for {
		select {
		case <-t.TokenBucket:
			continue
		case <-ticker.C:
			return nil
		}
	}
}

// start the timely fill token to the token bucket, if param maxRuntime below or equal to zero, the fill token will run forever
func (t *tokenBucket) FillToken(fillInterval time.Duration, maxRuntime time.Duration) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ticker := time.NewTicker(fillInterval)
		for {
			uid, _ := uuid.NewV4()
			<-ticker.C
			select {
			case t.TokenBucket <- uid:
			default:
			}

			fmt.Printf("bucket length: %d with %v\n", len(t.TokenBucket), time.Now())
		}
	}()
	if maxRuntime > 0 {
		ticker := time.NewTicker(maxRuntime)
		<-ticker.C
		return
	}
	wg.Wait()

}

// consumer to fetch a token from the token bucket, return a UUID and fetch result
func (t *tokenBucket) FetchToken() (string, bool) {
	var (
		taken bool
		token string
	)
	select {
	case uid := <-t.TokenBucket:
		fmt.Printf("fetch token:%s\n", uid)
		token = uid.String()
		taken = true
	default:
	}
	return token, taken
}
