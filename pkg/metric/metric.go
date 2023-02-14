package metric

// Metric is a sample interface
// Implementations of Metrics in metric package ar Counter, Gauge,
// PointGauge, RollingCounter and RollingGauge.
type Metric interface {
	//Add adds the given value to ther] counter
	Add(int64)
	//Value gets the current value.
	// If the metric's type is PointGauge, RollingCounter, RollingGauge,
	// it returns the sum value within the window.
	Value() int64
}

// Aggregation contains some common aggregation function.
// Each aggregation can compute summary statistics of window.
type Aggregation interface {
	Min() float64
	Max() float64
	Avg() float64
	Sum() float64
}

type CommonVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}
