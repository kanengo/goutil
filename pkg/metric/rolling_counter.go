package metric

import (
	"fmt"
	"time"
)

type RollingCounter interface {
	Metric
	Aggregation
	Reduce(func(BucketIterator) float64) float64
}

type rollingCounter struct {
	*RollingPolicy
}

type RollingCounterOpts struct {
	Size           int
	BucketDuration time.Duration
}

func NewRollingCounter(opts RollingCounterOpts) RollingCounter {
	window := NewWindow(WindowOpts{Size: opts.Size})
	policy := NewRollingPolicy(window, RollingPolicyOpts{BucketDuration: opts.BucketDuration})
	return &rollingCounter{policy}
}

func (r *rollingCounter) Add(val int64) {
	if val < 0 {
		panic(fmt.Errorf("metric counter:cannot decrease in value:%d", val))
	}
	r.RollingPolicy.Add(float64(val))
}

func (r *rollingCounter) Reduce(f func(BucketIterator) float64) float64 {
	return r.RollingPolicy.Reduce(f)
}

func (r *rollingCounter) Value() int64 {
	return int64(r.Sum())
}

func (r *rollingCounter) Min() float64 {
	return r.Reduce(Min)
}

func (r *rollingCounter) Max() float64 {
	return r.Reduce(Max)
}

func (r *rollingCounter) Avg() float64 {
	return r.Reduce(Avg)
}

func (r *rollingCounter) Sum() float64 {
	return r.Reduce(Sum)
}
