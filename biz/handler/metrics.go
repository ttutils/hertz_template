package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"net/url"
	"runtime"
)

// Create Prometheus metrics
var (
	memAlloc      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "go_memory_alloc_bytes", Help: "Current memory allocation in bytes."})
	numGoroutines = prometheus.NewGauge(prometheus.GaugeOpts{Name: "go_num_goroutines", Help: "Number of Goroutines."})
	totalAlloc    = prometheus.NewGauge(prometheus.GaugeOpts{Name: "go_memory_total_alloc_bytes", Help: "Total memory allocated in bytes."})
)

// Initialize Prometheus metrics
func init() {
	prometheus.MustRegister(memAlloc)
	prometheus.MustRegister(numGoroutines)
	prometheus.MustRegister(totalAlloc)
}

// Update Prometheus metrics
func updateMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memAlloc.Set(float64(m.Alloc))
	numGoroutines.Set(float64(runtime.NumGoroutine()))
	totalAlloc.Set(float64(m.TotalAlloc))
}

// Adapter for ResponseWriter
type responseWriterAdapter struct {
	resp *protocol.Response
}

func (r *responseWriterAdapter) Header() http.Header {
	headers := http.Header{}
	r.resp.Header.VisitAll(func(key, value []byte) {
		headers.Add(string(key), string(value))
	})
	return headers
}

func (r *responseWriterAdapter) Write(p []byte) (n int, err error) {
	r.resp.AppendBody(p)
	return len(p), nil
}

func (r *responseWriterAdapter) WriteHeader(statusCode int) {
	r.resp.SetStatusCode(statusCode)
}

// HertzRequestAdapter implements http.Request interface for Prometheus
type HertzRequestAdapter struct {
	req *protocol.Request
}

func (r *HertzRequestAdapter) Method() string {
	return string(r.req.Method())
}

func (r *HertzRequestAdapter) URL() *url.URL {
	u, _ := url.Parse(string(r.req.URI().Path()))
	return u
}

func (r *HertzRequestAdapter) Proto() string {
	return "HTTP/1.1"
}

func (r *HertzRequestAdapter) ProtoMajor() int {
	return 1
}

func (r *HertzRequestAdapter) ProtoMinor() int {
	return 1
}

func (r *HertzRequestAdapter) Header() http.Header {
	headers := http.Header{}
	r.req.Header.VisitAll(func(key, value []byte) {
		headers.Add(string(key), string(value))
	})
	return headers
}

func (r *HertzRequestAdapter) Body() io.ReadCloser {
	return io.NopCloser(nil)
}

func (r *HertzRequestAdapter) ContentLength() int64 {
	return int64(r.req.Header.ContentLength())
}

func (r *HertzRequestAdapter) Close() bool {
	return false
}

func (r *HertzRequestAdapter) Host() string {
	return string(r.req.Host())
}

func (r *HertzRequestAdapter) RequestURI() string {
	return string(r.req.RequestURI())
}

func (r *HertzRequestAdapter) RemoteAddr() string {
	return ""
}

func (r *HertzRequestAdapter) TLS() interface{} {
	return nil
}

// MetricsHandler Metrics Handler to expose Prometheus metrics
func MetricsHandler(ctx context.Context, c *app.RequestContext) {
	// Update metrics
	updateMetrics()

	// Create adapters for Hertz to http interfaces
	responseAdapter := &responseWriterAdapter{resp: &c.Response}

	// Create fake http.Request from Hertz request
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating metrics handler")
		return
	}

	// Serve the metrics to Prometheus
	promhttp.Handler().ServeHTTP(responseAdapter, req)
}
