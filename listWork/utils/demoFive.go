package utils

/*
题目：如何找出单链表中的倒数第K个元素
*/

//快慢指针法
/*
由于单链表只能从头到尾依次对各个节点进行访问，因此，需要找到链表的倒数第K个元素，也只能是从头到尾进行遍历查找
在查找的过程中，设置两个指针，让其中一个指针比另外一个指针先前移动K步骤，然后两个指针同时移动，循环直到先行的指针
为null为止，另外的一个指针就是需要找的位置
*/

func FindLastK(head *Node, kValue int) *Node {
	if head == nil || head.Next == nil {
		return head
	}
	slow := head.Next
	fast := head.Next
	i := 0
	for i = 0; i < kValue && fast != nil; i++ {
		fast = fast.Next
	}
	if i < kValue {
		return nil
	}
	for fast != nil {
		slow = slow.Next
		fast = fast.Next
	}
	return slow

}
