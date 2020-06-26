package DemoOne

/*
给定一个数组，编写一个函数，将所有的0移动到数组末尾，同时保持非0元素的相对顺序
*/

/*
解题思路
1 即选择一个数，比这个数小的放在左边，比数大的放在右边，这里可以改变判断条件
不等于0放左边，等于0放右边，
使用快慢指针：只要nums[i] != 0 交换nums[i]和nums[j]
283题
*/

func moveZeroes(nums []int) {
	if nums == nil {
		return
	}
	var j int
	for i := 0; i < len(nums); i++ {
		if nums[i] == 0 {
			continue
		}
		nums[i], nums[j] = nums[j], nums[i]
		j++
	}

}
