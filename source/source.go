package source

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Sourcer interface {
	Query(context.Context, *Request) ([]byte, error)
}

type SourceFunc func(context.Context, *Request) ([]byte, error)

// Source representes a data source
type Source struct {
	Name string
	Tag  string

	Client      *http.Client
	Request     *Request
	RequestFunc SourceFunc
	Logger      log.Logger
}

type Request struct {
	Method   string
	URL      *url.URL
	Headers  map[string]string
	Body     string
	Username string
	Password string
}

// NewSource creates a basic source
func New(name, tag, endpoint string) *Source {
	dc := defaultClient()
	dr := defautRequest(endpoint)
	return &Source{
		Name:        name,
		Tag:         tag,
		Client:      dc,
		RequestFunc: defaultSourceRequest(dc),
		Request:     dr,
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

func defautRequest(endpoint string) *Request {
	url, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	return &Request{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
	}
}

// Get performs and http request to a specified URL, with a timeout context
// of 5 Second by default
func (s *Source) Query(ctx context.Context, req *Request) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.RequestFunc(ctx, req)
}

func defaultSourceRequest(c *http.Client) SourceFunc {
	return func(ctx context.Context, request *Request) ([]byte, error) {
		req, err := http.NewRequest(request.Method, request.URL.String(), strings.NewReader(request.Body))
		if err != nil {
			return nil, err
		}

		if request.Password != "" && request.Username != "" {
			req.SetBasicAuth(request.Username, request.Password)
		}

		for k, v := range request.Headers {
			req.Header.Add(k, v)
		}

		resp, err := c.Do(req.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		return ioutil.ReadAll(resp.Body)
	}
}

// SetOpts adds custom logger, metrics
type SourceOption func(*Source)

func (s *Source) SetOptions(ops ...SourceOption) {
	for _, opFunc := range ops {
		opFunc(s)
	}
}

// SetBasicAuth returns a SourceOption which can be added on a source
func SetBasicAuth(username, password string) SourceOption {
	return func(source *Source) {
		source.Request.Username = username
		source.Request.Password = password
	}
}

func SetHeaders(headers map[string]string) SourceOption {
	return func(source *Source) {
		source.Request.Headers = headers
	}
}

// SetLogger returns a SourceOption which can be added on a source
func SetLogger(l log.Logger) SourceOption {
	return func(source *Source) {
		source.Logger = l
	}
}

// SetClient returns a SourceOption which can be added on a source
func SetClient(c *http.Client) SourceOption {
	return func(source *Source) {
		source.Client = c
	}
}

// SetSourceFunc returns a SourceOption which can be added on a source
func SetSourceFunc(reqFunc SourceFunc) SourceOption {
	return func(source *Source) {
		source.RequestFunc = reqFunc
	}
}
