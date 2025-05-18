package http

import "net/http"

type HandlerFunc func(*Context)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}
