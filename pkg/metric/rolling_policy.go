package metric

import (
	"sync"
	"time"
)

type RollingPolicy struct {
	size   int
	window *Window
	offset int

	bucketDuration time.Duration
	lastAppendTIme time.Time
	sync.RWMutex
}

type RollingPolicyOpts struct {
	BucketDuration time.Duration
}

func NewRollingPolicy(window *Window, opts RollingPolicyOpts) *RollingPolicy {
	return &RollingPolicy{
		size:           window.Size(),
		window:         window,
		offset:         0,
		bucketDuration: opts.BucketDuration,
		lastAppendTIme: time.Now(),
	}
}

func (r *RollingPolicy) timespan() int {
	v := int(time.Since(r.lastAppendTIme) / r.bucketDuration)
	if v > -1 {
		return v
	}

	return r.size
}

func (r *RollingPolicy) add(f func(offset int, val float64), val float64) {
	r.Lock()
	defer r.Unlock()
	timespan := r.timespan()
	if timespan > 0 {
		r.lastAppendTIme = r.lastAppendTIme.Add(time.Duration(timespan) * r.bucketDuration)
		offset := r.offset
		s := offset + 1
		if timespan > r.size {
			timespan = r.size
		}
		d, d1 := s+timespan, 0
		if d > r.size {
			d1 = d - r.size
			d = r.size
		}
		for i := s; i < d; i++ {
			r.window.ResetBucket(i)
			offset = i
		}
		for i := 0; i < d1; i++ {
			r.window.ResetBucket(i)
			offset = i
		}
		r.offset = offset
	}
	f(r.offset, val)
}

func (r *RollingPolicy) Append(val float64) {
	r.add(r.window.Append, val)
}

func (r *RollingPolicy) Add(val float64) {
	r.add(r.window.Add, val)
}

func (r *RollingPolicy) Reduce(f func(iterator BucketIterator) float64) float64 {
	r.RLock()
	defer r.RUnlock()
	timespan := r.timespan()
	var val float64
	if count := r.size - timespan; count > 0 {
		offset := r.offset + timespan + 1
		if offset >= r.size {
			offset = offset - r.size
		}
		val = f(r.window.Iterator(offset, count))
	}
	return val
}
