package serve

import (
	"log"
	"net/http"
	"strings"
)

type (
	HandlerFunc func(*Context)

	router struct {
		roots    map[string]*node       // method : *node
		handlers map[string]HandlerFunc // method-pattern: HandlerFunc
	}
)

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 分割 pattern, 返回 parts
// 最多只允许有一个 *
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)

	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = newNode()
	}

	parts := parsePattern(pattern)
	r.roots[method].insert(pattern, parts, 0)

	log.Printf("Add route %4s - %s", method, pattern)
	key := strings.Join([]string{method, pattern}, "-")
	r.handlers[key] = handler
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	searchParts := parsePattern(path) // 待匹配的 parts
	n := root.search(searchParts, 0)  // 找到匹配的 trie node

	if n == nil {
		return nil, nil
	}

	params := make(map[string]string)
	parts := parsePattern(n.pattern) // n 的 parts

	for index, part := range parts {
		// 将匹配结果与模式映射起来
		if part[0] == ':' {
			params[part[1:]] = searchParts[index]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}

	return n, params
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		c.Params = params
		key := strings.Join([]string{c.Method, n.pattern}, "-")
		handler := r.handlers[key]
		c.handlers = append(c.handlers, handler)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}

	c.Next()
}
