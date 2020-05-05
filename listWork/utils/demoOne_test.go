package utils

import "testing"

func TestHelloWorld(t *testing.T) {
	headValue := &CircleLink{}
	headValue.Insert(8)
	headValue.FmtList()
	headValue.Reverse()
	headValue.FmtList()

	t.Log("hello world")
}
