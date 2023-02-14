package metric

import (
	"fmt"
	"testing"
	"time"
)

func TestRollingPolicy(t *testing.T) {
	window := NewWindow(WindowOpts{Size: 10})
	rollingPolicy := NewRollingPolicy(window, RollingPolicyOpts{BucketDuration: time.Millisecond * 100})
	s := time.Now()
	for i := 0; i < 16; i++ {
		time.Sleep(time.Millisecond * 100)
		rollingPolicy.Add(float64(i + 1))
	}
	s2 := time.Since(s).Round(time.Millisecond * 100)
	fmt.Println("s2:", s2)
	rollingPolicy.Reduce(func(iterator BucketIterator) float64 {
		for iterator.Next() {
			b := iterator.Bucket()
			fmt.Println(b.Points)
		}
		return 0
	})
}
