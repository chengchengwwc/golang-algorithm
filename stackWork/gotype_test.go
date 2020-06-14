package stackWork

import "testing"

func TestHello(t *testing.T) {
	mm := []int{4, 5, 3, 2, 1}
	s := CreateStack(mm)
	SortStack(s)
	t.Log(s)
}
