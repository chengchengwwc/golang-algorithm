package utils

// 如何进行链表翻转

//1. 相邻元素位置翻转
/*
主要思路：通过调整节点指针域的指向来直接调换相邻的两个节点，如果单链表恰好有偶数个节点，那么只需要将
奇偶数节点对掉即可，如果链表有奇数个节点，则将除最后一个节点进行奇偶对掉即可。
*/
func Reverse(head *Node) {
	if head == nil || head.Next == nil {
		return
	}
	cur := head.Next
	pre := head
	var next *Node
	for cur != nil && cur.Next != nil {
		next = cur.Next.Next
		pre.Next = cur.Next
		cur.Next.Next = cur
		cur.Next = next
		cur = next
	}
}

//2 如何合并两个有序链表
func Merge(head1 *Node, head2 *Node) *Node {
	if head1 == nil || head1.Next == nil {
		return head1
	}
	if head2 == nil || head2.Next == nil {
		return head2
	}
	cur1 := head1.Next
	cur2 := head2.Next
	var head *Node
	var cur *Node
	if cur1.Data.(int) > cur2.Data.(int) {
		head = head2
		cur = cur2
		cur2 = cur2.Next
	} else {
		head = head1
		cur = cur1
		cur1 = cur1.Next
	}

	// 每次找到链表的最小结点的最小值对应的节点连接到合并链表的尾部
	for cur1 != nil && cur2 != nil {
		if cur1.Data.(int) < cur2.Data.(int) {
			cur.Next = cur1
			cur = cur1
			cur1 = cur1.Next
		} else {
			cur.Next = cur2
			cur = cur2
			cur2 = cur2.Next
		}
	}

	if cur1 != nil {
		cur.Next = cur1
	}

	if cur2 != nil {
		cur.Next = cur2
	}
	return head
}

//如何在只给定单链表中的某个节点指针的情况下删除该节点

func RemoveNode(head *Node, tmpValue int) bool {
	if head == nil || head.Next == nil {
		return false
	}
	for head != nil && head.Next != nil {
		if head.Data.(int) == tmpValue {
			head.Next = head.Next.Next
		}
		head = head.Next
	}
	return true
}

// 如何判断两个单链表是否有交叉
/*
1。hash方法： 如果两个链表有交叉，那么他们一定有公共的节点，由于节点的地址或是引用可以作为节点的唯一标示，因此可以
可以通过判断两个链表中的节点是否有相同的地址来判断
2. 首尾相接法：将两个链表首位相连，然后检测这个链表是否存在环，如果存在，则两个链表相交，而环的入口则为相交节点的
*/

func IsIntersect(head1 *Node, head2 *Node) *Node {
	if head1 == nil || head1.Next == nil || head2 == nil || head2.Next == nil {
		return nil
	}
	tmp1 := head1.Next
	tmp2 := head2.Next
	n1 := 0
	n2 := 0
	//遍历heaa1 找到尾节点，同时记录head1的长度
	for tmp1.Next != nil {
		tmp1 = tmp1.Next
		n1++
	}
	//遍历head2，找到尾节点，同时记录head2的长度
	for tmp2.Next != nil {
		tmp2 = tmp2.Next
		n2++
	}
	// head1和head2 是否有相同的尾节点
	if tmp1 == tmp2 {
		if n1 > n2 {
			for n1-n2 > 0 {
				head1 = head1.Next
				n1--
			}
		}
		if n2 > n1 {
			for n2-n1 > 0 {
				head2 = head2.Next
				n2--
			}
		}
		for head1 != head2 {
			head1 = head1.Next
			head2 = head2.Next
		}
		return head1
	}
	return nil
}

//
