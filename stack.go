package main

type stack struct {
    ar []interface{}
}

type Stack interface {
    push(x interface{})
    pop() interface{}
    size() int
}

func (s *stack) push(x interface{}) {
    s.ar = append(s.ar, x)
}

func (s *stack) pop() interface{} {
    if len(s.ar) == 0 {
        return nil
    }
    r := s.ar[len(s.ar)-1]
    s.ar = s.ar[:len(s.ar)-1]
    return r
}

func (s *stack) size() int {
    return len(s.ar)
}

func NewStack() Stack {
    return &stack{}
}