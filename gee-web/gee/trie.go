package gee

import (
    "fmt"
    "strings"
)

type node struct {
    pattern  string // 这是字典树记录的Data。记录是否有pattern，即有对应handler，相当于Data
    part     string // 这是构建字典树的索引，方便查找Data。
    children []*node
    isWild   bool
}

func (n *node) String() string {
    return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) matchChild(part string) *node {
    for _, child := range n.children {
        if child.part == part || child.isWild {
            return child
        }
    }
    return nil
}

func (n *node) matchChildren(part string) []*node {
    var children []*node
    for _, child := range n.children {
        if child.part == part || child.isWild {
            children = append(children, child)
        }
    }
    return children
}

func (n *node) insert(parts []string, pattern string) {
    if len(parts) == 0 {
        n.pattern = pattern
        return
    }
    part := parts[0]
    child := n.matchChild(part)
    if child == nil {
        child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
        n.children = append(n.children, child)
    }
    child.insert(parts[1:], pattern)
}

func (n *node) search(parts []string) *node {
    if len(parts) == 0 || strings.HasPrefix(n.part, "*") {
        if n.pattern == "" { // 即这里有Data
            return nil
        }
        return n
    }

    part := parts[0]
    children := n.matchChildren(part)
    for _, child := range children {
        // 按顺序搜索
        if target := child.search(parts[1:]); target != nil {
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
