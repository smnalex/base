package service

import (
	"log"
	"sync"

	"github.com/smnalex/base/source"
)

// Service represents the service
type Service struct {
	name    string
	tag     string
	Sources map[string]*source.Source
	sync.Mutex
	Logging log.Logger
}

// NewService creates the service
func NewService(name, tag string) *Service {
	return &Service{name: name, tag: tag}
}

// Logger returns the service logger
func (s *Service) Logger() log.Logger {
	return s.Logger()
}

// RegisterSources register all the sources for a service
// a source must be of type source.Sourcer(implements a Get method)
func (s *Service) RegisterSources(sources map[string]*source.Source) {
	s.Sources = sources
}

// GetSource returns a source from the service
func (s *Service) GetSource(name string) *source.Source {
	s.Lock()
	defer s.Unlock()
	return s.Sources[name]
}
