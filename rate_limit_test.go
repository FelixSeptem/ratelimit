package ratelimit

import (
	"fmt"
	"github.com/golang/go/src/math/rand"
	"sync"
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	tb := InitTokenBucket(3600)
	tb.Preheat(100)
	var wg sync.WaitGroup

	// producer
	go func() {
		wg.Add(1)
		tb.FillToken(time.Second * 1)
		wg.Done()
	}()

	// consumer
	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(time.Second * time.Duration(rand.Intn(7)))
			wg.Add(1)
			uid := tb.FetchToken()
			wg.Done()
			fmt.Printf("fetch token result:%v\n", uid)
		}
	}()
	wg.Wait()
	time.Sleep(time.Second * 1800)
}
