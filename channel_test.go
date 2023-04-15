package infinity_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Code-Hex/go-infinity-channel"
)

func TestNewChannel(t *testing.T) {
	ch := infinity.NewChannel[int]()
	if ch == nil {
		t.Fatal("NewChannel() should return a non-nil value")
	}
}

func TestInAndOut(t *testing.T) {
	ch := infinity.NewChannel[int]()

	go func() {
		for i := 0; i < 10; i++ {
			ch.In() <- i
		}
		ch.Close()
	}()

	var values []int
	for v := range ch.Out() {
		values = append(values, v)
	}

	for i, v := range values {
		if v != i {
			t.Errorf("Expected %d, got %d", i, v)
		}
	}
}

func TestLen(t *testing.T) {
	ch := infinity.NewChannel[int]()
	wantLen := 100
	go func() {
		for i := 0; i < wantLen; i++ {
			ch.In() <- i
		}
		ch.Close()
	}()

	time.Sleep(100 * time.Millisecond) // Give time for the goroutine to run

	if ch.Len() != wantLen {
		t.Fatalf("Expected length of 10, got %d", ch.Len())
	}
}

func TestConcurrency(t *testing.T) {
	ch := infinity.NewChannel[int]()
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ch.In() <- n
		}(i)
	}

	go func() {
		wg.Wait()
		ch.Close()
	}()

	var count int
	for range ch.Out() {
		count++
	}

	if count != 100 {
		t.Fatalf("Expected count of 10, got %d", count)
	}
}
