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

func InitTokenBucket(capacity int32) *tokenBucket {
	ch := make(chan uuid.UUID, capacity)
	return &tokenBucket{
		TokenBucket: ch,
	}
}

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

func (t *tokenBucket) FillToken(fillInterval time.Duration) {
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		ticker := time.NewTicker(fillInterval)
		for {
			uid, _ := uuid.NewV4()
			select {
			case <-ticker.C:
				select {
				case t.TokenBucket <- uid:
				default:
				}
			}
			fmt.Printf("bucket length: %d with %v\n", len(t.TokenBucket), time.Now())
		}
	}()
	time.Sleep(time.Second * 1)
	wg.Wait()

}

func (t *tokenBucket) FetchToken() bool {
	var taken bool
	select {
	case uid := <-t.TokenBucket:
		fmt.Printf("fetch token:%s\n", uid)
		taken = true
	default:
		taken = false
	}
	return taken
}
