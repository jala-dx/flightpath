package main

import (
    "fmt"
)

type node struct {
   val string 
   ar []*node
}

func Node(val string) *node {
    n := &node{}
    n.val = val
    return n
}

func Connect(n,x *node) {
    n.ar = append(n.ar,x)
}

func TopoSortHelper(n *node, s Stack, m map[*node]bool) {
    m[n] = true
    for _, p := range(n.ar) {
        if _, ok := m[p]; ok {
            continue
        }
        TopoSortHelper(p, s, m)
    }
    s.push(n)
}

func TopoSort(ar []*node) string {
    s := NewStack()
    m := map[*node]bool{}
    for _, n := range(ar) {
        if _, ok := m[n]; ok {
            continue;
        }
        TopoSortHelper(n, s, m)
    }
    flightPath := ""
    for s.size() != 0 {
        r := s.pop()
        n := r.(*node)
        fmt.Printf("%s\n",n.val) 
        flightPath = fmt.Sprintf("%s%s->", flightPath, n.val)
    }
    flightPath = flightPath[:len(flightPath)-2]
    return flightPath
}

func flights(inputMap map[string]string ) string {
    m := map[string]*node{}
    var nodes []*node
    var nodeObj *node
    for k, v := range(inputMap) {
        if _, ok := m[k]; !ok {
           nodeObj = Node(k)
           m[k] = nodeObj
           nodes = append(nodes, nodeObj) 
        }
        if _, ok := m[v]; !ok {
           nodeObj = Node(v)
           m[v] = nodeObj
           nodes = append(nodes, nodeObj)
        }
        Connect(m[k], m[v]) 
    }
    flightPath := TopoSort(nodes)
    return flightPath
}