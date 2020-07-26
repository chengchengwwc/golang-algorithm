package DemoFour


/*
N叉树

题目：590
*/

type Node struct {
	Val int
	Children []*Node
}

//迭代
func postorder(root *Node) []int {
	var res []int
	if root == nil {
		return []int{}
	}
	var stack = []*Node{root}
	for 0 < len(stack) {
		root = stack[len(stack) -1]
		res = append(res,root.Val)
		stack = stack[:len(stack)-1]
		l := len(root.Children)
		for i:=0;i<l;i++{
			stack = append(stack,root.Children[i])
		}
	}
	l := len(res) -1
	for i:=0;i<l/2+1;i++{
		res[i],res[l-i] = res[l-i],res[i]
	}
	return res
    
}