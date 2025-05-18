package http

import (
	"net/http"
	"strings"
)

type router struct {
	http.ServeMux
	middleware []MiddlewareFunc
}

func NewRouter() *router {
	return &router{}
}

func (r *router) GET(pattern string, handler HandlerFunc) {
	final := r.wrapWithMiddleware(handler, pattern, http.MethodGet)
	r.ServeMux.HandleFunc(cleanPattern(pattern), final)
}

func (r *router) POST(pattern string, handler HandlerFunc) {
	final := r.wrapWithMiddleware(handler, pattern, http.MethodPost)
	r.ServeMux.HandleFunc(cleanPattern(pattern), final)
}

func (r *router) PUT(pattern string, handler HandlerFunc) {
	final := r.wrapWithMiddleware(handler, pattern, http.MethodPut)
	r.ServeMux.HandleFunc(cleanPattern(pattern), final)
}

func (r *router) DELETE(pattern string, handler HandlerFunc) {
	final := r.wrapWithMiddleware(handler, pattern, http.MethodDelete)
	r.ServeMux.HandleFunc(cleanPattern(pattern), final)
}

func cleanPattern(pattern string) string {
	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = ""
		}
	}
	return "/" + strings.Join(parts, "/")
}
