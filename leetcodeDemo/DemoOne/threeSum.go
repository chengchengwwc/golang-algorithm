package DemoOne

import "sort"

/*
题目15
双指针左右下标，向中间推移
*/

func oldthreeSum(nums []int) [][]int {
	var result [][]int
	for i := 0; i < len(nums)-2; i++ {
		for j := i + 1; j < len(nums)-1; j++ {
			for k := j + 1; k < len(nums); k++ {
				if nums[i]+nums[j]+nums[k] == 0 {
					var tmp []int
					tmp = append(tmp, nums[i])
					tmp = append(tmp, nums[j])
					tmp = append(tmp, nums[k])
					result = append(result, tmp)
				}
			}
		}
	}
	return result
}

func goodthreeSum(nums []int) [][]int {
	var result [][]int
	tmpDict := make(map[int][]int)
	for i := 0; i < len(nums)-2; i++ {
		for j := 0; j < len(nums)-1; j++ {
			if tmpDict[nums[j]] != nil {
				newList := tmpDict[nums[j]]
				newList = append(newList, nums[j])
				result = append(result, newList)
			} else {
				tmpValue := 0 - nums[i] - nums[j]
				var tmplist []int
				tmplist = append(tmplist, nums[i])
				tmplist = append(tmplist, nums[j])
				tmpDict[tmpValue] = tmplist
			}
		}
	}
	return result
}

// 双指针法
func threeSum(nums []int) [][]int {
	if nums == nil {
		return nil
	}
	if len(nums) < 3 {
		return nil
	}

	sort.Ints(nums) // 排序
	var result [][]int
	var i int
	var L int
	var R int
	for i = 0; i < len(nums); i++ {
		if nums[i] > 0 {
			return result
		}
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		L = i + 1
		R = len(nums) - 1
		for L < R {
			var tmpList []int
			if nums[i]+nums[L]+nums[R] == 0 {
				tmpList = append(tmpList, nums[i])
				tmpList = append(tmpList, nums[L])
				tmpList = append(tmpList, nums[R])
				result = append(result, tmpList)
				for L < R && nums[L] == nums[L+1] {
					L = L + 1
				}
				for L < R && nums[R] == nums[R-1] {
					R = R - 1
				}
				L = L + 1
				R = R - 1
			} else if nums[i]+nums[R]+nums[L] > 0 {
				R = R - 1
			} else {
				L = L + 1
			}
		}
	}
	return result
}
