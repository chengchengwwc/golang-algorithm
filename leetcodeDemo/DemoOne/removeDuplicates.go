package DemoOne

/*
题目：26
快慢指针：删除重复数据
*/
func removeDuplicates(nums []int) int {
	n := len(nums)
	if n < 2 {
		return n
	}
	left, right := 0, 1
	for right < n {
		if nums[right] != nums[right-1] {
			nums[left] = nums[right]
			left++
		}
		right++
	}
	return left

}
