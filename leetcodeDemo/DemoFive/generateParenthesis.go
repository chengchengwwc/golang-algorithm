package DemoFive


/*
题目：22
递归主体：左括号小于n才能继续左括号，右括号小于左括号才能调教右括号
递归终止：左括号=右括号=n
*/


func generateParenthesis(n int) []string {
	var ouput []string
	_generate(0,0,n,"",&ouput)
	return ouput

}

func _generate(left int,right int,max int,s string, output *[]string){
	if left == right && left == max {
		*output = append(*output,s)
		return
	}

	// 递归主体
	if left < max{
		_generate(left+1, right, max, s + "(", output)
	}

	if right < left {
		_generate(left, right+1,  max, s + ")", output)
	}
}