package DemoOne

/*
题目：1

*/

func twoSum(nums []int, target int) []int {
	var reuslt []int
	m := make(map[int]int)
	for i, k := range nums {
		value, ok := m[target-k]
		if ok {
			reuslt = append(reuslt, value)
			reuslt = append(reuslt, i)
		}
		m[k] = i
	}
	return reuslt
}
