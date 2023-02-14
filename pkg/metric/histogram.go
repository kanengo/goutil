package metric

import "github.com/prometheus/client_golang/prometheus"

type HistogramVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	Buckets   []float64
}

type HistogramVec interface {
	Observe(v float64, labels ...string)
}

type promHistogramVec struct {
	histogram *prometheus.HistogramVec
}

func NewHistogramVec(opts *HistogramVecOpts) HistogramVec {
	if opts == nil {
		return nil
	}

	vec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: opts.Namespace,
		Subsystem: opts.Subsystem,
		Name:      opts.Name,
		Help:      opts.Help,
		Buckets:   opts.Buckets,
	}, opts.Labels)

	prometheus.MustRegister(vec)
	return &promHistogramVec{histogram: vec}
}

func (p *promHistogramVec) Observe(v float64, labels ...string) {
	p.histogram.WithLabelValues(labels...).Observe(v)
}
