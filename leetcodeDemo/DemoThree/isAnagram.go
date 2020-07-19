package DemoThree

/*
题目 242
*/



func isAnagram(s string, t string) bool {
	if len(s) != len(t){
		return false
	}
	wordstable_s := [26]int{}
	wordstable_t := [26]int{}
	for i:=0;i<len(s);i++{
		index := s[i] - 'a'
		wordstable_s[index] ++
	}
	for i:=0;i<len(t);i++{
		index := t[i] - 'a'
		wordstable_t[index] ++
	}
	return wordstable_s == wordstable_t
}