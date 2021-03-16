package combination

import (
	"log"
	"testing"
)

type stringArray struct {
	self     []string
	refArray *[]string
}

func (sa stringArray) String() string {
	retString := "["
	comma := ""
	for _, v := range sa.self {
		retString += comma + v
		comma = ", "
	}
	return retString + "]"
}

func (sa stringArray) CopyElement(src int, dst int) {
	sa.self[dst] = (*sa.refArray)[src]
}

func TestWithStringArray(t *testing.T) {
	refArray := []string{"bob", "fred", "steve"}
	da := stringArray{
		make([]string, len(refArray)-1),
		&refArray,
	}
	todoChan := NewCombChannelLen(len(refArray), 2)
	ga := NewGenericWorker(todoChan)
	for err := ga.Next(da); err == nil; err = ga.Next(da) {
		log.Println(da)
	}
}

func TestWithStrings(t *testing.T) {
	refArray := []string{"bob", "fred", "steve"}
	dstArray := make([]string, len(refArray))
	copyFunc := func(i, j int) {
		dstArray[i] = refArray[j]
	}
	gc := NewGeneric(len(refArray), 2, copyFunc)
	cnt := 0
	for err := gc.Next(); err == nil; err = gc.Next() {
		//dstArray will not have stuff in it
		t.Log(dstArray)
		cnt++
	}
	if cnt != 3 {
		t.Error("Cnt should be 3, was:", cnt)
	}
}
