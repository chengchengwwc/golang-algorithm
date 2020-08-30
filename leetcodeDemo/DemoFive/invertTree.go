package DemoFive

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}
// 前序遍历
func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	root.Left,root.Right  = root.Right,root.Left
	invertTree(root.Right)
	invertTree(root.Left)
	return root
}