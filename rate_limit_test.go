package ratelimit

import (
	"fmt"
	"math/rand"
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
		tb.FillToken(time.Millisecond * 100)
		wg.Done()
	}()

	for i := 0; i < 3; i++ {
		// start multiple consumer
		go func(n int) {
			for i := 0; i < 1000; i++ {
				time.Sleep(time.Second * time.Duration(rand.Intn(7)))
				wg.Add(1)
				uid := tb.FetchToken()
				wg.Done()
				fmt.Printf("[Consumer %d]fetch token result:%v\n", n, uid)
				if i%10 == 0 {
					tb.Flush(time.Second)
				}
			}
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second * 10)
}
