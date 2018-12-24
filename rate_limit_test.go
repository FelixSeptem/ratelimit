package ratelimit

import (
	"testing"
	"time"
)

func TestInitTokenBucket(t *testing.T) {
	tb := InitTokenBucket(1024)
	if cap(tb.TokenBucket) != 1024 || len(tb.TokenBucket) != 0 {
		t.Errorf("expect len = 0 and cap = 100, got len = %d cap = %d", len(tb.TokenBucket), cap(tb.TokenBucket))
	}
}

func TestTokenBucket_FetchToken(t *testing.T) {
	tb := InitTokenBucket(1024)
	tb.FillToken(time.Millisecond*20, time.Millisecond*200)
	time.Sleep(time.Millisecond * 20)
	uid, fetch := tb.FetchToken()
	if !fetch {
		t.Error("not fetch token")
	}
	if len(uid) == 0 {
		t.Error("got empty token")
	}
}

func TestTokenBucket_Preheat(t *testing.T) {
	tb := InitTokenBucket(1024)
	err := tb.Preheat(100)
	if err != nil {
		t.Fatal(err)
	}
	if len(tb.TokenBucket) != 100 {
		t.Errorf("expect len = 100 got len = %d", len(tb.TokenBucket))
	}
}

func TestTokenBucket_Flush(t *testing.T) {
	tb := InitTokenBucket(1024)
	err := tb.Preheat(100)
	if err != nil {
		t.Fatal(err)
	}
	err = tb.Flush(time.Millisecond * 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(tb.TokenBucket) != 0 {
		t.Errorf("expect len = 0 got len = %d", len(tb.TokenBucket))
	}
}
