package DemoOne

/*
题目：88
*/

func merge(nums1 []int, m int, nums2 []int, n int) {
	maxNum := m + n
	left, right := m-1, n-1
	i := 1
	for left >= 0 && right >= 0 {
		if nums1[left] < nums2[right] {
			nums1[maxNum-i] = nums2[right]
			i++
			right--
		} else {
			nums1[maxNum-i] = nums1[left]
			i++
			left--
		}
	}
	if left == -1 && right >= 0 {
		for i := right; i >= 0; i-- {
			nums1[i] = nums2[i]
		}
	}
}
