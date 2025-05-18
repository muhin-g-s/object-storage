package http

import (
	"net/http"
	"strings"
)

type MiddlewareFunc func(*Context, HandlerFunc)

func (r *Router) Use(mw ...MiddlewareFunc) {
	r.middleware = append(r.middleware, mw...)
}

func (r *Router) wrapWithMiddleware(handler HandlerFunc, pattern string, method string) http.HandlerFunc {
	finalHandler := handler
	for i := len(r.middleware) - 1; i >= 0; i-- {
		mw := r.middleware[i]
		next := finalHandler

		finalHandler = func(ctx *Context) {
			mw(ctx, next)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		params := parseParams(pattern, r.URL.Path)

		ctx := &Context{
			Writer:  w,
			Request: r,
			Params:  params,
		}

		finalHandler(ctx)
	}
}

func parseParams(pattern, path string) map[string]string {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	params := map[string]string{}

	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], ":") {
			key := strings.TrimPrefix(patternParts[i], ":")
			if i < len(pathParts) {
				params[key] = pathParts[i]
			}
		}
	}
	return params
}
