package source

// Do not change!!! change only when you find something stupid ;)

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smnalex/base/log"
)

type Sourcer interface {
	Get(context.Context, *url.Values) ([]byte, error)
}

// Sourcer an interface to be implemented by a data source
type SourceFunc func(context.Context, *url.Values) ([]byte, error)

// Source representes a data source
type Source struct {
	Client *http.Client

	name     string
	tag      string
	Endpoint string
	Username string
	Password string
	Logger   log.Logger
	Metrics  prometheus.Counter
	Tracing  interface{}
	SF       SourceFunc
}

// Get performs and http request to a specified endpoint, with a timeout context
// of 5 Second by default
func (s *Source) Get(ctx context.Context, q *url.Values) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := s.SF(ctx, q)

	// Handle common error cases probably use switch type
	if err != nil {
		s.Logger.Debug("Failed Request", map[string]interface{}{"err": err, "context": ctx, "urlVal": q})
		s.Metrics.Inc()
		return nil, err
	}
	return res, nil
}

func defaultSourceRequest(c *http.Client, e string) SourceFunc {
	return func(ctx context.Context, q *url.Values) ([]byte, error) {
		req, err := http.NewRequest("GET", e, nil)
		if err != nil {
			return nil, err
		}

		req = req.WithContext(ctx)
		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Invalid request %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}
}

// NewSource creates a basic source
func NewSource(name, tag, endpoint string) *Source {
	df := defaultClient()
	return &Source{
		name:     name,
		tag:      tag,
		Endpoint: endpoint,
		Client:   df,
		SF:       defaultSourceRequest(df, endpoint),
	}
}

// defaultClient returns http.Client, use SetClient to add a custom one
func defaultClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
		},
	}
}

// SetAuth used by a request when basic auth is needed
func (s *Source) SetAuth(username, password string) {
	s.Username = username
	s.Password = password
}

// SetOpts adds custom logger, metrics
type SourceOption func(*Source)

func (s *Source) SetOptions(ops ...SourceOption) {
	for _, e := range ops {
		e(s)
	}
}

// SetLogger returns a SourceOption which can be added on the source
func SetLogger(l log.Logger) SourceOption {
	return func(source *Source) {
		source.Logger = l
	}
}

// SetMetric returns a SourceOption which can be added on the source
func SetMetrics(m prometheus.Counter) SourceOption {
	return func(source *Source) {
		source.Metrics = m
	}
}

// SetTracer returns a SourceOption which can be added on the source
func SetTracer(t log.Logger) SourceOption {
	return func(source *Source) {
		source.Tracing = t
	}
}

// SetClient returns a SourceOption which can be added on the source
func SetClient(c *http.Client) SourceOption {
	return func(source *Source) {
		source.Client = c
	}
}
