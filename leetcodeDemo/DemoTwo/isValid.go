package DemoTwo

/*
题目 20
最近相关性-- 栈

*/

func isValid(s string) bool {
	m := map[byte]byte{
		'(': ')',
		'[': ']',
		'{': '}',
	}
	lps := make([]byte, 0, len(s)/2)
	for i := 0; i < len(s); i++ {
		p := s[i]
		if _, ok := m[p]; ok {
			// 实现进栈
			lps = append(lps, p)
		} else {
			//实现出栈
			if len(lps) > 0 && m[lps[len(lps)-1]] == p {
				lps = lps[:len(lps)-1]
			} else {
				return false
			}
		}
	}
	if len(lps) == 0 {
		return true
	}
	return false
}
