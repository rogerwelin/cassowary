package client

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/schollz/progressbar"
	"go.uber.org/ratelimit"
)

type LoadTest func(c *Cassowary, outPutChan chan<- durationMetrics, g *QueryGroup)

// Validator(statusCode int, respSize int64, resp []byte, err error) (failed ool, statusCode string)
type Validator func(int, int64, []byte, error) (bool, string)

type Query struct {
	Name           string
	Method         string
	URL            string
	DataType       string
	Data           []byte // Body
	RequestHeaders [][2]string
	Validator      Validator // Custom Validator function
}

type QueryGroup struct {
	Name string

	ConcurrencyLevel int
	Delay            time.Duration
	Requests         int

	l ratelimit.Limiter

	FileMode    bool
	URLPaths    []string
	URLIterator Iterator

	Method        string
	Data          []byte
	RequestHeader []string

	loadTest LoadTest // Custom load test function

	workerChan chan *Query
}

// Cassowary is the main struct with bootstraps the load test
type Cassowary struct {
	IsTLS                 bool
	BaseURL               string
	ExportMetrics         bool
	ExportMetricsFile     string
	PromExport            bool
	Cloudwatch            bool
	Histogram             bool
	Boxplot               bool
	StatFile              string
	TLSConfig             *tls.Config
	PromURL               string
	DisableTerminalOutput bool
	DisableKeepAlive      bool
	Client                *http.Client
	Bar                   *progressbar.ProgressBar
	Timeout               int

	Duration time.Duration

	// Old syntax for default group
	ConcurrencyLevel int
	Requests         int
	RequestHeader    []string
	URLPaths         []string

	// New syntax with query groups
	Groups []QueryGroup

	wgStart sync.WaitGroup
	wgStop  sync.WaitGroup
}

// ResultMetrics are the aggregated metrics after the load test
type ResultMetrics struct {
	Name              string         `json:"name"`
	BaseURL           string         `json:"base_url"`
	TotalRequests     int            `json:"total_requests"`
	FailedRequests    int            `json:"failed_requests"`
	RespSuccess       map[string]int `json:"responses_success"`
	RespFailed        map[string]int `json:"responses_failed"`
	RequestsPerSecond float64        `json:"requests_per_second"`
	DNSMedian         float64        `json:"dns_median"`
	ElapsedStats      stats          `json:"elapsed"`
	TCPStats          stats          `json:"tcp_connect"`
	ProcessingStats   stats          `json:"server_processing"`
	ContentStats      stats          `json:"content_transfer"`
	BodySize          stats          `json:"body_size"`
	RespSize          stats          `json:"resp_size"`
}

type stats struct {
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	P95    float64 `json:"P95"`
	P99    float64 `json:"P99"`
}
