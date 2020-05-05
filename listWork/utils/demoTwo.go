package utils

import "fmt"

/*
如何从无序链表中删除重复项
题目： 给定一个没有排序的链表，去掉其重复项，并保留原来的顺序
*/

//顺序删除

//顺序删除
func (self *CircleLink) RemoveDump() {
	if self.head == nil || self.head.Next == nil {
		return
	}
	outerCur := self.head.Next //外层循环，用于链表的第一个节点
	var innerCur *Node         // 内层循环，用于遍历outerCur后面的节点
	var innerPre *Node         // innerCur的前驱节点
	for ; outerCur != nil; outerCur = outerCur.Next {
		fmt.Println(innerPre)
		for innerCur, innerPre = outerCur.Next, outerCur; innerCur != nil; {
			if outerCur.Data == innerCur.Data {
				innerCur = innerCur.Next
			} else {
				innerPre = innerCur
				innerCur = innerCur.Next
			}
		}
	}
}
