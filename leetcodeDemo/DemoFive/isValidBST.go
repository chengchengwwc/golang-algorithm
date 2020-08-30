package DemoFive

import "math"

/*
题目：98
*/

func isValidBST(root *TreeNode) bool {
	return validBST(root,math.MinInt64,math.MaxInt64)

}

func validBST(root *TreeNode,min int,max int) bool {
	if root == nil{
		return true
	}
	if root.Val <= min || root.Val >= max{
		return false
	}
	return validBST(root.Left,min,root.Val) && validBST(root.Right,root.Val,max)

}