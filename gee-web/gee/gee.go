package gee

import (
    "html/template"
    "log"
    "net/http"
    "path"
    "strings"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
    prefix      string
    middlewares []HandlerFunc
    parent      *RouterGroup
    engine      *Engine
}

type Engine struct {
    *RouterGroup
    router        *router
    groups        []*RouterGroup
    htmlTemplates *template.Template
    funcMap       template.FuncMap
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var middlewares []HandlerFunc
    for _, group := range e.groups {
        if strings.HasPrefix(r.URL.Path, group.prefix) {
            middlewares = append(middlewares, group.middlewares...)
        }
    }
    c := newContext(w, r)
    c.handlers = middlewares
    c.engine = e
    e.router.handle(c)
}

func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
    g.addRoute("GET", pattern, handler)
}
func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
    g.addRoute("POST", pattern, handler)
}

func (g *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
    pattern := path.Join(g.prefix, comp)
    log.Printf("addRoute method: %s pattern: %s", method, pattern)
    g.engine.router.addRoute(method, pattern, handler)
}

func New() *Engine {
    e := &Engine{}
    e.RouterGroup = &RouterGroup{engine: e}
    e.router = newRouter()
    e.groups = []*RouterGroup{e.RouterGroup}
    return e
}

func Default() *Engine {
    e := New()
    e.Use(Logger(), Recovery())
    return e
}

func (e *Engine) Run(addr string) error {
    return http.ListenAndServe(addr, e)
}

func (g *RouterGroup) Group(prefix string, middlewares ...HandlerFunc) *RouterGroup {
    engine := g.engine
    group := &RouterGroup{
        prefix:      path.Join(g.prefix, prefix),
        middlewares: middlewares,
        parent:      g,
        engine:      engine,
    }
    engine.groups = append(engine.groups, group)
    return group
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
    g.middlewares = append(g.middlewares, middlewares...)
}

func (g *RouterGroup) createStaticHandler(relPath string, fs http.FileSystem) HandlerFunc {
    absPath := path.Join(g.prefix, relPath)
    // 去掉request里面的relPath
    fileServer := http.StripPrefix(absPath, http.FileServer(fs))
    return func(c *Context) {
        file := c.Param("filepath")
        // param里的filepath只是来检查一下是否存在
        if _, err := fs.Open(file); err != nil {
            c.Status(http.StatusNotFound)
            return
        }
        // 实际serve还是用的request里的path
        fileServer.ServeHTTP(c.Writer, c.Req)
    }
}

func (g *RouterGroup) Static(relPath string, root string) {
    // url到relPath的，去root下面找文件
    pattern := path.Join(relPath, "/*filepath")
    handler := g.createStaticHandler(relPath, http.Dir(root))
    g.GET(pattern, handler)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
    e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
    e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}
