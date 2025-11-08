package geecache

import (
    "fmt"
    "net/http"
    "strings"
)

const defaultBasePath = "/_geecache/"

type HTTPPool struct {
    self     string // 记录自己的地址
    basePath string // 节点之间通信的地址前缀
}

func NewHTTPPool(self string) *HTTPPool {
    return &HTTPPool{self: self, basePath: defaultBasePath}
}

func (p *HTTPPool) Log(format string, args ...any) {
    fmt.Printf("[server %s] %s", p.self, fmt.Sprintf(format, args...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if !strings.HasPrefix(r.URL.Path, p.basePath) {
        panic("HTTPPool serving unexpected path: " + r.URL.Path)
    }
    p.Log("%s %s", r.Method, r.URL.Path)
    parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
    if len(parts) != 2 {
        http.Error(w, "bad request", http.StatusBadRequest)
    }
    groupName := parts[0]
    key := parts[1]
    group := GetGroup(groupName)
    if group == nil {
        http.Error(w, "no such group:"+groupName, http.StatusNotFound)
        return
    }
    view, err := group.Get(key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Write(view.ByteSlice())
    return
}
