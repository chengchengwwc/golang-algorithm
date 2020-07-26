package DemoFour


/*
二叉树：前序输出
题目：104



*/


//前序输出
func preorderTraversal(root *TreeNode) []int {
	var res []int
	var stack []*TreeNode
	if root == nil{
		return []int{}
	}
	for 0 < len(stack) || root != nil{
		for root != nil{
			res = append(res,root.Val)
			stack  = append(stack,root.Right)
			root = root.Left
		}
		index := len(stack) -1
		root = stack[index]
		stack = stack[:index]
	}
	return res
}

func newpreorderTraversal(root *TreeNode) []int{
	var max *TreeNode
	var res []int
	for root != nil{
		if root.Left == nil{
			res = append(res,root.Val)
			root = root.Right
		} else {
			max = root.Left
			for max.Right != nil{
				max = max.Right
			}
			root.Right,max.Right = root.Left,max.Left
			root.Left = nil
		}
	}
	return res
}