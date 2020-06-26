package DemoOne

/*
141
快慢指针法
*/

func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}
	slow := head
	fast := head.Next
	for slow != fast {
		if fast == nil || fast.Next == nil {
			return false
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	return true
}

// hashMap方法
func hashMapCycle(head *ListNode) bool {
	hash := make(map[*ListNode]int)
	for head != nil {
		if _, ok := hash[head]; ok {
			return true
		}
		hash[head] = head.Val
		head = head.Next
	}
	return false
}
