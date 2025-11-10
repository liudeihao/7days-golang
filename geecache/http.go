package geecache

import (
    "fmt"
    "geecache/geecache/consistenthash"
    "io"
    "net/http"
    "net/url"
    "strings"
    "sync"
)

const (
    defaultBasePath = "/_geecache/"
    defaultReplicas = 50
)

type HTTPPool struct {
    self        string // 记录自己的地址
    basePath    string // 节点之间通信的地址前缀
    mu          sync.Mutex
    peers       *consistenthash.Map
    httpGetters map[string]*httpGetter
}

type httpGetter struct {
    baseURL string
}

func (g *httpGetter) Get(group string, key string) ([]byte, error) {
    u := fmt.Sprintf(
        "%v%v/%v",
        g.baseURL,
        url.QueryEscape(group),
        url.QueryEscape(key),
    )
    res, err := http.Get(u)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("server returned status code %d", res.StatusCode)
    }

    bytes, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, fmt.Errorf("reading response body: %v", err)
    }
    return bytes, nil

}

var _ PeerGetter = (*httpGetter)(nil)

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

func (p *HTTPPool) Set(peers ...string) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.peers = consistenthash.New(defaultReplicas, nil)
    p.peers.Add(peers...)
    p.httpGetters = make(map[string]*httpGetter, len(peers))
    for _, peer := range peers {
        p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
    }
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
    p.mu.Lock()
    defer p.mu.Unlock()
    if peer := p.peers.Get(key); peer != "" && peer != p.self {
        p.Log("Pick peer %s", peer)
        return p.httpGetters[peer], true
    }
    return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)
