package actuator

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Set of linked build time variables for providing relevant information for the application
var (
	BuildStamp       string
	GitCommitAuthor  string
	GitCommitID      string
	GitCommitTime    string
	GitPrimaryBranch string
	GitURL           string
	HostName         string
	Username         string
	Version          string
)

// HealthCheckInfo is the information for health check
type HealthCheckInfo struct {
	Key         string    `json:"-"`
	IsMandatory bool      `json:"-"`
	Success     bool      `json:"success"`
	Error       string    `json:"error,omitempty"`
	LastCheckTs time.Time `json:"-"`
}

var healthCheckInfoRWMutex sync.RWMutex
var healthCheckInfo map[string]*HealthCheckInfo

func getBasicInfo(config *Config) map[string]interface{} {
	version := strings.TrimSpace(config.Version)
	if version == "" {
		version = Version
	}
	return map[string]interface{}{
		applicationKey: map[string]string{
			envKey:     config.Env,
			nameKey:    config.Name,
			versionKey: version,
		},
		gitKey: map[string]string{
			buildStampKey:       BuildStamp,
			gitCommitAuthorKey:  GitCommitAuthor,
			gitCommitIDKey:      GitCommitID,
			gitCommitTimeKey:    GitCommitTime,
			gitPrimaryBranchKey: GitPrimaryBranch,
			gitURLKey:           GitURL,
			hostNameKey:         HostName,
			usernameKey:         Username,
		},
	}
}

func getHealthCheckInfoForPing(h *HealthChecker) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.Timeout)
	defer cancel()
	return h.Ping(ctx)
}

func getHealthCheckInfoForEndpoint(h *HealthChecker) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.Timeout)
	defer cancel()
	req, err := http.NewRequest(http.MethodGet, h.Endpoint, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("%s failed with status %d", h.Key, res.StatusCode)
}

func getHealthCheckInfoForHealthChecker(h *HealthChecker, ch chan *HealthCheckInfo) {
	var lastCheckTs time.Time
	var i *HealthCheckInfo

	healthCheckInfoRWMutex.RLock()
	i = healthCheckInfo[h.Key]
	if i != nil {
		lastCheckTs = i.LastCheckTs
	}
	healthCheckInfoRWMutex.RUnlock()

	// in case of previous failure or no cache exist or last check beyond duration
	// refresh the health check
	if i == nil || !i.Success || h.CacheDuration == 0 || time.Now().Sub(lastCheckTs) > h.CacheDuration {
		var err error
		if h.Ping != nil {
			err = getHealthCheckInfoForPing(h)
		} else {
			err = getHealthCheckInfoForEndpoint(h)
		}
		healthCheckInfoRWMutex.Lock()
		i = &HealthCheckInfo{
			Key:         h.Key,
			IsMandatory: h.IsMandatory,
			Success:     err == nil,
			LastCheckTs: time.Now(),
		}
		if err != nil {
			i.Error = err.Error()
		}
		healthCheckInfo[h.Key] = i
		healthCheckInfoRWMutex.Unlock()
	}

	ch <- i
}

func getHealthCheckInfo(config *Config) (result map[string]*HealthCheckInfo, success bool) {
	if len(config.HealthCheckers) == 0 {
		return nil, true
	}
	ch := make(chan *HealthCheckInfo, len(config.HealthCheckers))
	defer close(ch)
	result = make(map[string]*HealthCheckInfo)
	success = true
	for _, h := range config.HealthCheckers {
		go getHealthCheckInfoForHealthChecker(h, ch)
	}
	for range config.HealthCheckers {
		h := <-ch
		result[h.Key] = h
		if !h.Success && h.IsMandatory {
			success = false
		}
	}
	return
}

// getInfoHandler is used to get the handler function for the info endpoint
func getInfoHandler(config *Config) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Add(contentTypeHeader, applicationJSONContentType)
		info := getBasicInfo(config)
		healthInfo, success := getHealthCheckInfo(config)
		if len(healthInfo) > 0 {
			info["health"] = healthInfo
		}
		infoB, _ := encodeJSON(info)
		if !success {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = writer.Write(infoB)
	}
}
