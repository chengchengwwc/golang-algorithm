package DemoFour

/*
题目：94
二叉树中序遍历


*/



type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

//递归
func inorderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	rest := append(inorderTraversal(root.Left),root.Val)
	rest = append(rest,inorderTraversal(root.Right)...)
	return rest

}

// 迭代
func newinorderTraversal(root *TreeNode) []int{
	return inorderIterate(root)

}

func inorderIterate(root *TreeNode) []int {
	if root == nil {
		return []int {}
	}
	var res []int
	var stack []*TreeNode
	for 0 < len(stack) || root != nil {
		for root != nil{
			stack = append(stack,root)
			root = root.Left
		}
		index := len(stack) - 1
		res = append(res,stack[index].Val)
		root = stack[index].Right
		stack = stack[:index]
	}
	return res

}




