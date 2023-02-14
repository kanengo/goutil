package metric

import (
	"time"
)

type RollingGauge interface {
	Metric
	Aggregation
	Reduce(func(BucketIterator) float64) float64
}

type rollingGauge struct {
	*RollingPolicy
}

type RollingGaugeOpts struct {
	Size           int
	BucketDuration time.Duration
}

func NewRollingGauge(opts RollingGaugeOpts) RollingGauge {
	window := NewWindow(WindowOpts{Size: opts.Size})
	policy := NewRollingPolicy(window, RollingPolicyOpts{BucketDuration: opts.BucketDuration})
	return &rollingGauge{policy}
}

func (r *rollingGauge) Add(val int64) {
	r.RollingPolicy.Add(float64(val))
}

func (r *rollingGauge) Reduce(f func(BucketIterator) float64) float64 {
	return r.RollingPolicy.Reduce(f)
}

func (r *rollingGauge) Value() int64 {
	return int64(r.Sum())
}

func (r *rollingGauge) Min() float64 {
	return r.Reduce(Min)
}

func (r *rollingGauge) Max() float64 {
	return r.Reduce(Max)
}

func (r *rollingGauge) Avg() float64 {
	return r.Reduce(Avg)
}

func (r *rollingGauge) Sum() float64 {
	return r.Reduce(Sum)
}
