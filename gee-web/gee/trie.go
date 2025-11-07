package gee

import (
    "fmt"
    "strings"
)

type node struct {
    pattern  string  // 待匹配的路由，如 /p/:lang
    part     string  // 路由中的一部分， 例如:lang
    children []*node // 子节点，例如[doc, tutorial, intro]
    isWild   bool    // 是否精确匹配，即part含有:或*
}

func (n *node) String() string {
    return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
    for _, child := range n.children {
        if child.part == part || child.isWild {
            return child
        }
    }
    return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
    var children []*node
    for _, child := range n.children {
        if child.part == part || child.isWild {
            children = append(children, child)
        }
    }
    return children
}

func (n *node) insert(pattern string, parts []string, height int) {
    if len(parts) == height {
        n.pattern = pattern
        return
    }
    part := parts[height]
    child := n.matchChild(part)
    if child == nil {
        child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
        n.children = append(n.children, child)
    }
    child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
    if len(parts) == height || strings.HasPrefix(n.part, "*") {
        if n.pattern == "" { // pattern为空的时候也是匹配失败：没有对应的路由
            return nil
        }
        return n
    }

    part := parts[height]
    children := n.matchChildren(part)
    for _, child := range children {
        // 按顺序搜索
        if target := child.search(parts, height+1); target != nil {
            return target
        }
    }
    return nil
}

func (n *node) travel(list *[]*node) {
    if n.pattern != "" {
        *list = append(*list, n)
    }
    for _, child := range n.children {
        child.travel(list)
    }
}
