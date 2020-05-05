package utils

/*
题目：对链表进行重新排序
*/

func findMiddleNode(head *Node) *Node {
	if head == nil || head.Next == nil {
		return head
	}
	fast := head //遍历链表，向前走两步
	slow := head // 遍历链表，向前走一步
	slowPre := head
	for fast != nil && fast.Next != nil {
		slowPre = slow
		slow = slow.Next
		fast = fast.Next.Next
	}
	slowPre.Next = nil
	return slow
}

// 对不带头节点的列表进行翻转
func reverse(head *Node) *Node {
	if head != nil && head.Next == nil {
		return head
	}
	var pre *Node
	var next *Node
	for head != nil {
		next = head.Next
		head.Next = pre
		pre = head
		head = next
	}
	return pre
}

func Record(head *Node) {
	if head == nil || head.Next == nil {
		return
	}

	curl := head.Next
	mid := findMiddleNode(head.Next)
	curl2 := reverse(mid)
	var tmp *Node
	for curl.Next != nil {
		tmp = curl.Next
		curl.Next = curl2
		curl = tmp
		tmp = curl2.Next
		curl2.Next = curl
		curl2 = tmp
	}
	curl.Next = curl2
}
