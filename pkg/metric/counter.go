package metric

import "github.com/prometheus/client_golang/prometheus"

type CounterVec interface {
	Inc(labels ...string)
	Add(v float64, labels ...string)
}

type CounterVecOpts CommonVecOpts

type promCounterVec struct {
	counter *prometheus.CounterVec
}

func (p *promCounterVec) Inc(labels ...string) {
	p.counter.WithLabelValues(labels...).Inc()
}

func (p *promCounterVec) Add(v float64, labels ...string) {
	p.counter.WithLabelValues(labels...).Add(v)
}

func NewCounterVec(opts *CounterVecOpts) CounterVec {
	if opts == nil {
		return nil
	}

	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: opts.Namespace,
		Subsystem: opts.Subsystem,
		Name:      opts.Name,
		Help:      opts.Help,
	}, opts.Labels)
	prometheus.MustRegister(vec)
	return &promCounterVec{counter: vec}
}
