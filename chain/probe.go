package chain

import "time"

// ProbeType names the probing method for node liveness checks.
type ProbeType string

const (
	ProbeTypeTCP ProbeType = "tcp"
	ProbeTypeHTTP ProbeType = "http"
	ProbeTypeCmd  ProbeType = "cmd"
)

// ProbeConfig holds the configuration for a node-level liveness probe.
// It is embedded in node config and drives the per-node probe goroutine.
type ProbeConfig struct {
	Type           ProbeType
	Addr           string
	Interval       time.Duration
	Timeout        time.Duration
	HTTPPath       string
	HTTPHost       string
	HTTPHeaders    map[string]string
	ExpectedStatus int
	Command        string
}

// ProbeResult records the outcome of a single probe attempt.
type ProbeResult struct {
	Success   bool
	Latency   time.Duration
	Error     string
	Timestamp time.Time
}

// ProbeResultReader is implemented by types that expose a probe result,
// enabling generic latency-aware strategies and filters.
type ProbeResultReader interface {
	ProbeResult() *ProbeResult
}
