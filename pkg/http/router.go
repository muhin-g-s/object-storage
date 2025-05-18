package http

import (
	"net/http"
	"strings"
)

type Router struct {
	http.ServeMux
	middleware []MiddlewareFunc
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) GET(pattern string, handler HandlerFunc) {
	final := r.wrapWithMiddleware(handler, pattern, http.MethodGet)
	r.ServeMux.HandleFunc(cleanPattern(pattern), final)
}

func (r *Router) POST(pattern string, handler HandlerFunc) {
	final := r.wrapWithMiddleware(handler, pattern, http.MethodPost)
	r.ServeMux.HandleFunc(cleanPattern(pattern), final)
}

func (r *Router) PUT(pattern string, handler HandlerFunc) {
	final := r.wrapWithMiddleware(handler, pattern, http.MethodPut)
	r.ServeMux.HandleFunc(cleanPattern(pattern), final)
}

func (r *Router) DELETE(pattern string, handler HandlerFunc) {
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
