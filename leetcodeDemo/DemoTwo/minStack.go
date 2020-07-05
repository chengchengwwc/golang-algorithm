package DemoTwo

import "math"

/*
设计一个支持 push ，pop ，top 操作，并能在常数时间内检索到最小元素的栈
题目：155

*/

type MinStack struct {
	stack    []int
	minStack []int
}

/** initialize your data structure here. */
func Constructor() MinStack {
	return MinStack{
		stack:    []int{},
		minStack: []int{math.MaxInt64},
	}

}

func (this *MinStack) Push(x int) {
	this.stack = append(this.stack, x)
	top := this.minStack[len(this.minStack)-1]
	this.minStack = append(this.minStack, min(x, top))

}

func min(x int, top int) int {
	if x > top {
		return top
	} else {
		return x
	}

}

func (this *MinStack) Pop() {
	this.stack = this.stack[:len(this.stack)-1]
	this.minStack = this.minStack[:len(this.minStack)-1]
}

func (this *MinStack) Top() int {
	return this.stack[len(this.stack)-1]

}

func (this *MinStack) GetMin() int {
	return this.minStack[len(this.minStack)-1]

}
