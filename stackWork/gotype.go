package stackWork

import "errors"

// 栈的数据结构:先进后出
// 用数组实现

type SliceStack struct {
	arr       []int
	stackSize int
}

type Node struct {
	Data interface{}
	Next *Node
}

func (p *SliceStack) IsEmpty() bool {
	return p.stackSize == 0
}

func (p *SliceStack) Size() int {
	return p.stackSize
}

func (p *SliceStack) Top() int {
	if p.IsEmpty() {
		panic(errors.New("bad"))
	}
	return p.arr[p.stackSize-1]
}

func (p *SliceStack) Pop() int {
	if p.stackSize > 0 {
		p.stackSize--
		ret := p.arr[p.stackSize]
		p.arr = p.arr[:p.stackSize]
		return ret
	}
	panic(errors.New(":Bad"))
}
func (p *SliceStack) Push(t int) {
	p.arr = append(p.arr, t)
	p.stackSize += 1
}

// 链表实现
type LinkedStack struct {
	head *Node
}

func (p *LinkedStack) IsEmpty() bool {
	return p.head.Next == nil
}
func (p *LinkedStack) Size() int {
	size := 0
	node := p.head.Next
	for node != nil {
		node = node.Next
		size++
	}
	return size
}

func (p *LinkedStack) Push(e int) {
	node := &Node{
		Data: e,
		Next: p.head.Next,
	}
	p.head.Next = node
}

func (p *LinkedStack) Pop() int {
	tmp := p.head.Next
	if tmp != nil {
		p.head.Next = tmp.Next
		return tmp.Data.(int)
	}
	panic(errors.New("bad"))
}

func (p *LinkedStack) Top() int {
	if p.head.Next != nil {
		return p.head.Next.Data.(int)
	}
	panic(errors.New("bad"))
}

// 栈中数据翻转
func moveBottom(s *SliceStack) {
	if s.IsEmpty() {
		return
	}
	top1 := s.Pop()
	if !s.IsEmpty() {
		//递归处理
		moveBottom(s)
		top2 := s.Pop()
		s.Push(top1)
		s.Push(top2)
	} else {
		s.Push(top1)
	}
}

func moveBottomToTopSort(s *SliceStack) {
	if s.IsEmpty() {
		return
	}
	top1 := s.Pop()
	if !s.IsEmpty() {
		moveBottomToTopSort(s)
		top2 := s.Top()
		if top1 > top2 {
			s.Pop()
			s.Push(top1)
			s.Push(top2)
			return
		}
	}
	s.Push(top1)
}

func SortStack(s *SliceStack) {
	if s.IsEmpty() {
		return
	}
	moveBottomToTopSort(s)
	top := s.Pop()
	SortStack(s)
	s.Push(top)
}

func ReverseStack(s *SliceStack) {
	if s.IsEmpty() {
		return
	}
	moveBottom(s)
	top := s.Pop()
	ReverseStack(s)
	s.Push(top)
}

func CreateStack(targetList []int) *SliceStack {
	stack := &SliceStack{}
	for _, v := range targetList {
		stack.Push(v)
	}
	return stack
}
