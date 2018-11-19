package algorithm

import (
	"container/list"
)

//https://segmentfault.com/a/1190000017052768

type Node struct {
	Value int
	Left *Node
	Right *Node
}

type Stack struct {
	list *list.List
}

func NewStack() *Stack {
	list := list.New()
	return &Stack{list}
}

func(stack *Stack) Push(value interface{}){
	stack.list.PushBack(value)
}

func (stack *Stack)Pop() interface{} {
	if e := stack.list.Back(); e!= nil {
		stack.list.Remove(e)
		return e.Value
	}
	return nil
}

func (stack *Stack) Len() int {
	return stack.list.Len()
}

func (stack *Stack) Empty() bool {
	return stack.Len() == 0
}