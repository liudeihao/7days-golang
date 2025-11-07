package gee

import (
    "net/http"
    "strings"
)

type router struct {
    roots    map[string]*node
    handlers map[string]HandlerFunc
}

func newRouter() *router {
    return &router{
        roots:    make(map[string]*node),
        handlers: make(map[string]HandlerFunc),
    }
}

func parsePattern(pattern string) []string {
    ss := strings.Split(pattern, "/")
    var parts []string
    for _, s := range ss {
        if s != "" {
            parts = append(parts, s)
            if s[0] == '*' {
                break
            }
        }
    }

    return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
    // 每种方法有对应的字典树
    if _, ok := r.roots[method]; !ok {
        r.roots[method] = &node{}
    }
    parts := parsePattern(pattern)
    r.roots[method].insert(pattern, parts, 0)

    key := method + " - " + pattern
    r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
    pathParts := parsePattern(path)
    params := make(map[string]string)
    root, ok := r.roots[method]
    if !ok {
        return nil, nil
    }
    n := root.search(pathParts, 0)
    if n == nil {
        return nil, nil
    }
    // pattern中的param与path中的对应起来
    patternParts := parsePattern(n.pattern)
    for i, part := range patternParts {
        if part[0] == ':' {
            params[part[1:]] = pathParts[i]
        } else if len(part) > 1 && part[0] == '*' {
            params[part[1:]] = strings.Join(pathParts[i:], "/")
            break
        }
    }
    return n, params
}

func (r *router) getRoutes(method string) []*node {
    root, ok := r.roots[method]
    if !ok {
        return nil
    }
    nodes := make([]*node, 0)
    root.travel(&nodes)
    return nodes
}

func (r *router) handle(c *Context) {
    n, params := r.getRoute(c.Method, c.Path)
    if n == nil {
        c.handlers = append(c.handlers, func(c *Context) {
            c.String(http.StatusNotFound, "404 page not found: %s", c.Path)
        })
        return
    }
    key := c.Method + " - " + n.pattern
    c.handlers = append(c.handlers, r.handlers[key])
    c.Params = params
    c.Next()
}
