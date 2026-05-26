package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start1 := time.Now()
	<-or1(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start1))

	start2 := time.Now()
	<-or2(
		sig(1*time.Hour),
		sig(2*time.Minute),
		sig(2*time.Second),
		sig(1*time.Hour),
		sig(2*time.Second),
	)
	fmt.Printf("done after %v\n", time.Since(start2))
}

func or1(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	out := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan interface{}) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case <-ch:
				cancel()
			}
		}(ch)
	}

	go func() {
		defer cancel() //nolint:govet
		wg.Wait()
		close(out)
	}()

	return out
}

func or2(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	out := make(chan interface{})

	go func() {
		defer close(out)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or2(append(channels[3:], out)...):
			}
		}
	}()

	return out
}
