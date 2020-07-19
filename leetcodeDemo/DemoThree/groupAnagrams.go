package DemoThree

import (
	"sort"
)

/*
题目：49

通过排序法进行处理，将字符串转换为自定义的byte类型，然后排序，再转换为字符串
*/
type bytes []byte

func (b bytes) Len() int {
	return len(b)
}

func (b bytes)Less(i,j int) bool{
	return b[i] < b[j]
}

func (b bytes) Swap(i,j int) {
	b[i],b[j] = b[j],b[i]
}



func groupAnagrams(strs []string) [][]string { 
	res := [][]string{}
	m := make(map[string] int)
	for _,str := range strs {
		kBytes := bytes(str)
		sort.Sort(kBytes)
		k := string(kBytes)
		if idx,ok := m[k];!ok{
			m[k] = len(res)
			res = append(res,[]string{str})
		} else {
			res[idx] = append(res[idx],str)
		}
	}
	return res
}