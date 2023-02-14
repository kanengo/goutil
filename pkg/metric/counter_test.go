package metric

import (
	"testing"
)

func TestCounter(t *testing.T) {
	//reqTotal := NewCounterVec(&CounterVecOpts{
	//	Namespace: "Test",
	//	Subsystem: "request",
	//	Name:      "total",
	//	Help:      "fuck you bitch",
	//	Labels:    []string{"service", "method"},
	//})
	//reqDuartion := NewHistogramVec(&HistogramVecOpts{
	//	Namespace: "Test",
	//	Subsystem: "request",
	//	Name:      "duration_seconds",
	//	Help:      "mother fucker",
	//	Labels:    []string{"service", "method"},
	//	Buckets:   []float64{.005, .01, .025, .1, .25, .5, 1},
	//})
	//services := []string{"user", "auth", "cms", "vip"}
	//methods := []string{"List", "Update", "Delete", "Add"}
	//go func() {
	//	for {
	//		service := services[rand.Intn(len(services))]
	//		method := methods[rand.Intn(len(methods))]
	//		Inc(service, method)
	//		start := time.Now()
	//		time.Sleep(time.Millisecond * time.Duration(100*rand.Intn(3)+100))
	//		Observe(time.Since(start).Seconds(), service, method)
	//	}
	//}()
	//http.Handle("/metrics", promhttp.Handler())
	//err := http.ListenAndServe("0.0.0.0:9012", nil)
	//if err != nil {
	//	zaplog.Fatal("listen http prometheus failed", zap.Error(err))
	//}
}
