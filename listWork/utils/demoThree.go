package utils

import "log"

/*
题目：如何计算两个单链表所代表的数之和
*/

// 链表相加法
/*
该方法需要对两个链表都进行遍历，因此时间复杂度为O(n)，空间复杂度也是O(n)

*/

func AddLink(oneLink *CircleLink, twoLink *CircleLink) *Node {
	if oneLink == nil || oneLink.head.Next == nil {
		log.Fatal(oneLink)
		return nil
	}
	if twoLink == nil || twoLink.head.Next == nil {
		log.Fatal(twoLink)
		return nil
	}

	c := 0   // 需要补位的数据
	sum := 0 // 相加之和
	p1 := oneLink.head.Next
	p2 := twoLink.head.Next
	resultHead := &Node{} //相加后链表头节点
	p := resultHead
	for p1 != nil && p2 != nil {
		p.Next = &Node{}
		sum = p1.Data.(int) + p2.Data.(int) + c
		p.Next.Data = sum % 10
		c = sum / 10
		p = p.Next
		p1 = p1.Next
		p2 = p2.Next
	}
	// 链表one比链表two长，那么后面只需要考虑链表one的值
	if p2 == nil {
		for p1 != nil {
			p.Next = &Node{}
			sum = p1.Data.(int) + c
			p.Next.Data = sum % 10
			c = sum / 10
			p = p.Next
			p1 = p1.Next
		}
	}

	//链表two比链表one长，那么后面只需要考虑链表two的值
	if p1 == nil {
		for p2 != nil {
			p.Next = &Node{}
			sum = p2.Data.(int) + c
			p.Next.Data = sum % 10
			c = sum / 10
			p = p.Next
			p2 = p2.Next
		}
	}

	if c == 1 {
		p.Next = &Node{}
		p.Next.Data = c
	}
	return resultHead

}
