package main

import (
	"fmt"
	"github.com/FelixSeptem/ratelimit"
	"math/rand"
	"sync"
	"time"
)

func main() {
	tb := ratelimit.InitTokenBucket(3600)
	if err := tb.Preheat(100); err != nil {
		fmt.Println(err.Error())
	}
	var wg sync.WaitGroup

	// producer
	wg.Add(1)
	go func() {
		tb.FillToken(time.Millisecond*100, time.Second*3)
		wg.Done()
	}()

	for i := 0; i < 3; i++ {
		// start multiple consumer
		go func(n int) {
			for i := 0; i < 1000; i++ {
				time.Sleep(time.Second * time.Duration(rand.Intn(7)))
				wg.Add(1)
				uid, _ := tb.FetchToken()
				wg.Done()
				fmt.Printf("[Consumer %d]fetch token result:%v\n", n, uid)
				if i%10 == 0 {
					tb.Flush(time.Second)
				}
			}
		}(i)
	}
	wg.Wait()
}
