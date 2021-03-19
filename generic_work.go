package combination

import "math/rand"

type GenericWorker struct {
	todo <-chan ToDo
}

type CombinationGenericInterface interface {
	CopyElement(src int, dst int)
}

func NewGenericWorker(todo <-chan ToDo) *GenericWorker {
	v := new(GenericWorker)
	v.todo = todo
	return v
}
func (g GenericWorker) Next(dstArray CombinationGenericInterface) error {
	populateFrom, ok := <-g.todo
	if !ok {
		return ErrItterationComplete
	}
	for dst, src := range populateFrom {
		dstArray.CopyElement(int(src), int(dst))
	}
	return nil
}

type GenericCombination struct {
	resultChan <-chan ToDo
	copyFunc   func(int, int)
}

func NewGenericFromChan(ch <-chan ToDo, copyFunc func(int, int)) *GenericCombination {
	v := new(GenericCombination)
	v.resultChan = ch
	v.copyFunc = copyFunc
	return v
}
func NewGeneric(src, dst int, copyFunc func(int, int)) *GenericCombination {
	return NewGenericFromChan(NewCombChannelLen(src, dst), copyFunc)
}
func RandomToDoSource(cnt, srcLen, dstLen int) <-chan ToDo {
	todoChan := make(chan ToDo)
	go func() {
		for i := 0; i < cnt; i++ {
			tmp := rand.Perm(srcLen)
			tmpP := make(ToDo, dstLen)
			for j := 0; j < dstLen; j++ {
				tmpP[j] = Position(tmp[j])
			}
			todoChan <- tmpP
		}
		close(todoChan)
	}()
	return todoChan
}
func (gc GenericCombination) Next() error {
	if gc.copyFunc == nil {
		return ErrMissingCopyFunc
	}
	resArray, ok := <-gc.resultChan
	if !ok {
		return ErrItterationComplete
	}
	for i, v := range resArray {
		gc.copyFunc(i, int(v))
	}
	return nil
}

func (gc GenericCombination) NextSkipN(n int) error {
	if gc.copyFunc == nil {
		return ErrMissingCopyFunc
	}
	var resArray ToDo
	var ok bool
	for i := 0; i < n; i++ {
		resArray, ok = <-gc.resultChan
		if !ok {
			return ErrItterationComplete
		}
	}
	for i, v := range resArray {
		gc.copyFunc(i, int(v))
	}
	return nil
}
