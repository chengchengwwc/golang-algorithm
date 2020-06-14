package stackWork

import "errors"

// 如何实现队列

// 方法一：数组实现

type SliceQueue struct {
	arr  []int
	font int // 队列头
	rear int // 队列尾
}

func (p *SliceQueue) IsEmpty() bool {
	return p.font == p.rear
}

func (p *SliceQueue) Size() int {
	return p.rear - p.font
}

func (p *SliceQueue) GetFront() int {
	if p.IsEmpty() {
		panic(errors.New("bad"))
	}
	return p.arr[p.font]
}

func (p *SliceQueue) GetBack() int {
	if p.IsEmpty() {
		panic(errors.New("bad"))
	}
	return p.arr[p.rear-1]
}

func (p *SliceQueue) DeQueue() {
	if p.rear > p.font {
		p.rear--
		p.arr = p.arr[1:]
	} else {
		panic(errors.New("bad"))
	}

}

func (p *SliceQueue) AddQueue(item int) {
	p.arr = append(p.arr, item)
	p.rear++
}

// 方法二： 链表实现

type LinkedQueue struct {
	head *Node
	end  *Node
}

func (p *LinkedQueue) IsEmpty() bool {
	return p.head == nil
}

func (p *LinkedQueue) Size() int {
	size := 0
	node := p.head
	for node != nil {
		node = node.Next
		size++
	}
	return size
}

func (p *LinkedQueue) EnQueue(e int) {
	node := &Node{Data: e}
	if p.head == nil {
		p.head = node
		p.end = node
	} else {
		p.end.Next = node
		p.end = node
	}
}

func (p *LinkedQueue) DeQueue() {
	p.head = p.head.Next
	if p.head == nil {
		p.end = nil
	}
}
