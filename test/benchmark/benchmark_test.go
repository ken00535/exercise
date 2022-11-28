package benchmark

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRequestParallel(b *testing.B) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	var wip int32
	var cnt int32
	var wg sync.WaitGroup
	wg.Add(1)
	firstCh := make(chan int, 1)
	periodCh := make(chan time.Duration)
	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		var first time.Duration
		var periods []time.Duration
		start := time.Now()
		for {
			select {
			case <-ticker.C:
				fmt.Printf("ts: %v, wip: %d, cnt: %d \n",
					time.Now().Sub(start),
					atomic.LoadInt32(&wip),
					atomic.LoadInt32(&cnt),
				)
			case <-firstCh:
				first = time.Now().Sub(start)
			case period := <-periodCh:
				periods = append(periods, period)
			case <-ctx.Done():
				total := time.Now().Sub(start)
				var allPeriod time.Duration
				for _, p := range periods {
					allPeriod += p
				}
				average := allPeriod / time.Duration(b.N)
				fmt.Printf("first: %v, average: %v, total: %v\n", first, average, total)
				wg.Done()
				return
			}
		}
	}()
	b.SetParallelism(10)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()
			atomic.AddInt32(&wip, 1)
			res, _ := client.Get("http://192.168.0.32/XfsdKLfa")
			atomic.AddInt32(&wip, -1)
			if atomic.AddInt32(&cnt, 1) == 1 {
				firstCh <- 1
			}
			period := time.Now().Sub(start)
			periodCh <- period
			assert.Equal(b, http.StatusFound, res.StatusCode)
		}
	})
	cancel()
	wg.Wait()
}
