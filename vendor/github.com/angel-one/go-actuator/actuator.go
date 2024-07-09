package actuator

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Endpoints enumeration
const (
	Ping = iota
	Info
	Metrics
	ThreadDump
	Health
)

// DefaultEndpoints is the list of endpoints enabled by default
var DefaultEndpoints = []int{Info, Ping, Health}

// AllEndpoints is the list of endpoints supported
var AllEndpoints = []int{Info, Ping, Metrics, ThreadDump, Health}

// HealthChecker is used as the config for doing health check
type HealthChecker struct {
	Key           string
	Ping          func(ctx context.Context) error
	Endpoint      string
	CacheDuration time.Duration
	IsMandatory   bool
	Timeout       time.Duration
}

// Config is the set of configurable parameters for the actuator setup
type Config struct {
	Endpoints      []int
	Env            string
	Name           string
	Port           int
	Version        string
	HealthCheckers []*HealthChecker
}

func (config *Config) validate() {
	for _, endpoint := range config.Endpoints {
		if !isValidEndpoint(endpoint) {
			panic(fmt.Errorf("invalid endpoint %d provided", endpoint))
		}
	}
	hcm := make(map[string]bool)
	for _, h := range config.HealthCheckers {
		if hcm[h.Key] {
			panic(fmt.Errorf("key repeated %s", h.Key))
		}
		hcm[h.Key] = true
		if h.Ping != nil && h.Endpoint != "" {
			panic(fmt.Errorf("both endpoint and ping provided for key %s", h.Key))
		}
		if h.Ping == nil && h.Endpoint == "" {
			panic(fmt.Errorf("either endpoint or ping should be provided for ket %s", h.Key))
		}
	}
}

// Default is used to fill the default configs in case of any missing ones
func (config *Config) setDefaults() {
	if config.Endpoints == nil {
		config.Endpoints = DefaultEndpoints
	}
	for _, h := range config.HealthCheckers {
		if h.Timeout == 0 {
			h.Timeout = 5 * time.Second // 5 second default timeout
		}
	}
}

// GetActuatorHandler is used to get the handler function for the actuator endpoints
// This single handler is sufficient for handling all the endpoints.
func GetActuatorHandler(config *Config) http.HandlerFunc {
	if config == nil {
		config = &Config{}
	}
	handleConfigs(config)
	handlerMap := getHandlerMap(config)
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			// method not allowed for the requested resource
			sendStringResponse(writer, http.StatusMethodNotAllowed, methodNotAllowedError)
			return
		}
		endpoint := fmt.Sprintf("/%s", getLastStringAfterDelimiter(request.URL.Path, slash))
		if handler, ok := handlerMap[endpoint]; ok {
			handler(writer, request)
			return
		}
		// incorrect endpoint
		// or endpoint not enabled
		sendStringResponse(writer, http.StatusNotFound, notFoundError)
	}
}

func handleConfigs(config *Config) {
	config.validate()
	config.setDefaults()
	healthCheckInfo = make(map[string]*HealthCheckInfo)
}

func getHandlerMap(config *Config) map[string]http.HandlerFunc {
	handlerMap := make(map[string]http.HandlerFunc, len(config.Endpoints))
	for _, e := range config.Endpoints {
		// now one by one add the handler of each endpoint
		switch e {
		case Health:
			handlerMap[healthEndpoint] = getInfoHandler(config)
		case Info:
			handlerMap[infoEndpoint] = getInfoHandler(config)
		case Metrics:
			handlerMap[metricsEndpoint] = handleMetrics
		case Ping:
			handlerMap[pingEndpoint] = handlePing
		case ThreadDump:
			handlerMap[threadDumpEndpoint] = handleThreadDump
		}
	}
	return handlerMap
}
