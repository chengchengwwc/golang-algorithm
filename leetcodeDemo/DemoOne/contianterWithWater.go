package DemoOne

/*
给你 n 个非负整数 a1，a2，...，an，每个数代表坐标中的一个点 (i, ai) 。
在坐标内画 n 条垂直线，垂直线 i 的两个端点分别为 (i, ai) 和 (i, 0)。找出其中的两条线，
使得它们与 x 轴共同构成的容器可以容纳最多的水。

*/

/*
思路：从两边进行收敛，计算出当中的最大值

*/

func maxArea(height []int) int {
	if height == nil {
		return 0
	}
	if len(height) == 2 {
		if height[0] > height[1] {
			return height[1]
		} else if height[0] == height[1] {
			return height[0]
		} else {
			return height[0]
		}
	}

	var max int
	i := 0
	j := len(height) - 1
	for i != j {
		if height[i] < height[j] {
			temp := height[i] * (j - i)
			if temp > max {
				max = temp
			}
			i++
		} else if height[i] > height[j] {
			temp := height[j] * (j - i)
			if temp > max {
				max = temp
			}
			j--
		} else {
			temp := height[i] * (j - i)
			if temp > max {
				max = temp
			}
			i++

		}
	}
	return max

}
