package utils



// 如何检查一个较大的单链表是否有环
/*
单链表有环是说单链表中某个节点的next域指向的是链表中的在它之前的某一个节点，这样
在链表的尾部形成一个环形结构

方法：快慢指针遍历法
定义两个指针fast 和slow,两者的初始值都指向链头表，指针slow每次前进一步，指针fast前进两步。两个指针同时
向前移动，如果快指针等于慢指针，就证明这个链表是带环的单向链表，否则，证明这个链表是不带环的
*/

// 判断单链表是否有环
func IsLoop(head *Node) *Node{
	if head == nil || head.Next == nil{
		return head
	}
	slow := head.Next
	fast := head.Next
	for fast != nil && fast.Next != nil{
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return slow
		}
	}
	return nil
}

// 找出环的入口
func FindLoopNode(head *Node,meetNode *Node) *Node{
	first := head.Next
	second := meetNode
	for first != second{
		first = first.Next
		second = second.Next
	}
	return first

}

