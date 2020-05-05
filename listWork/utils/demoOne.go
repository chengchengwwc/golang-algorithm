package utils

/**
如何实现链表逆序
*/

import (
	"fmt"
)

type object interface{}

type Node struct {
	Data object
	Next *Node
}

type CircleLink struct {
	head *Node
}

func (self *CircleLink) Insert(tempLenth int) {
	for i := 0; i < tempLenth; i++ {
		temp := self.head
		newNode := &Node{}
		newNode.Data = i
		if temp == nil {
			self.head = newNode
		} else {
			for {
				if temp.Next == nil {
					break
				}
				temp = temp.Next
			}
			temp.Next = newNode
		}

	}
}

func (self *CircleLink) FmtList() {
	temp := self.head
	for {
		fmt.Println(temp.Data)
		if temp.Next == nil {
			break
		}
		temp = temp.Next
	}
}

//就地逆序
/*
性能分析：
时间复杂度为O(n)，其中n为链表长度，但是需要常数个额外的当量来保存当前节点的前驱和后驱，因此空间复杂度为O(1)

*/
func (self *CircleLink) Reverse() {
	if self.head == nil || self.head.Next == nil {
		return
	}
	var pre *Node // 定义前驱节点
	var cur *Node // 定义当前节点
	next := self.head.Next
	for next != nil {
		cur = next.Next
		next.Next = pre
		pre = next
		next = cur
	}
	self.head.Next = pre
}

//插入法实现逆序

func (this *CircleLink) InsertReverse() {
	if this.head.Next == nil || this.head == nil {
		return
	}
	var cur *Node
	var next *Node
	cur = this.head.Next
	for cur != nil {
		next = cur.Next
		cur.Next = this.head.Next
		this.head.Next = cur
		cur = next
	}
}
