package DemoFour




/*
N叉树前序遍历
题目：589
*/

func preorder(root *Node) []int {
	if root == nil {
		return []int{}
	}
	var res = []int{}
	var statck = []*Node{root}
	for len(statck) > 0{
		for root != nil{
			res = append(res,root.Val)
			if len(root.Children) == 0 {
				break
			}
			for i := len(root.Children)-1;i>0;i-- {
				statck = append(statck,root.Children[i])
			}
			root = root.Children[0]
		}
		root = statck[len(statck) -1]
		statck = statck[:len(statck)-1]
	}
	return res
}