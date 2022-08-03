// -*- tab-width: 2 -*-

// Package counters enables 1 line creation of stats to track your program flow; you get summaries every minute
// implemented using channels now but switching to sync based on these tests maybe; will keep both implementations
package counters

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

var cbRan int32 = 0

func metricReporterCB(metrics []MetricReport) {
	atomic.StoreInt32(&cbRan, 1)
	for x := range metrics {
		m := metrics[x]
		fmt.Println("CB: ", m.Name, m.Delta)
	}
}

func BenchmarkCounter(b *testing.B) {
	InitCounters()
	SetLogInterval(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Incr("num_of_things")
	}
}

func BenchmarkSyncIncr(b *testing.B) {
	InitCounters()
	SetLogInterval(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IncrSync("num_of_things")
	}
}

func TestCounter(t *testing.T) {

	InitCounters()
	SetLogInterval(1)
	SetMetricReporter(metricReporterCB)
	AddMetaCounter("availability", "good", "bad", RatioTotal)
	for i := 0; i < 1000; i++ {
		go func() {
			Incr("num_of_things")
			Incr("a_num_of_things")
		}()
	}
	time.Sleep(1100 * time.Millisecond)
	IncrDelta("good", 97)
	IncrDelta("bad", 3)
	Decr("num_of_things")
	c := atomic.LoadInt32(&cbRan)
	if c != 1 {
		fmt.Println("Callback did not run")
		t.Fail()
	}
	time.Sleep(1100 * time.Millisecond)

}
