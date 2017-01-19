package mware

import "net/http"

type handler func(http.Handler) http.Handler

type mware struct {
	handlers []handler
}

func New(hs ...handler) *mware {
	return &mware{append(([]handler)(nil), hs...)}
}

func (mw *mware) Append(hs ...handler) {
	mw.handlers = append(mw.handlers, hs...)
}

func (mw *mware) Run(h http.Handler) http.Handler {
	if h == nil {
		return http.DefaultServeMux
	}

	for i := range mw.handlers {
		h = mw.handlers[len(mw.handlers)-1-i](h)
	}
	return h
}

func (mw *mware) RunFunc(h http.HandlerFunc) http.Handler {
	return mw.Run(http.HandlerFunc(h))
}
