package gateway

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) error {
	if handler == nil {
		return errors.New("handler nil error.")
	}
	log.Printf("Route %4s - %s", method, pattern)
	key := strings.Join([]string{method, pattern}, "-")
	r.handlers[key] = handler
	return nil
}

func (r *router) handle(c *Context) error {
	if c == nil {
		return errors.New("context nil error.")
	}
	key := strings.Join([]string{c.Method, c.Path}, "-")
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
	return nil
}
