package combination

type GenericWorker struct {
	todo <-chan []int
}

type CombinationGenericInterface interface {
	// Len() int
	CopyElement(src int, dst int)
}

func NewGenericWorker(todo <-chan []int) *GenericWorker {
	v := new(GenericWorker)
	v.todo = todo
	return v
}
func (g GenericWorker) Next(dstArray CombinationGenericInterface) error {
	populateFrom, ok := <-g.todo
	if !ok {
		return ItterationCompleteError
	}
	for dst, src := range populateFrom {
		dstArray.CopyElement(src, dst)
	}
	return nil
}

type GenericCombination struct {
	resultChan <-chan []int
	copyFunc   func(int, int)
}

func NewGeneric(src, dst int, copyFunc func(int, int)) *GenericCombination {
	v := new(GenericCombination)
	v.resultChan = NewCombChannelLen(src, dst)
	v.copyFunc = copyFunc
	return v
}

func (gc GenericCombination) Next() error {
	if gc.copyFunc == nil {
		return MissingCopyFuncError
	}
	resArray, ok := <-gc.resultChan
	if !ok {
		return ItterationCompleteError
	}
	for i, v := range resArray {
		gc.copyFunc(i, v)
	}
	return nil
}
