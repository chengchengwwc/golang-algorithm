package DemoOne

/*
1. 暴力解题
2。 找最近重复子问题
3 70题
*/

func climbStairs(n int) int {
	if n <= 2 {
		return n
	}
	n1 := 1
	n2 := 2
	n3 := 3
	for i := 3; i < n+1; i++ {
		n3 = n1 + n2
		n1 = n2
		n2 = n3
	}
	return n3

}
