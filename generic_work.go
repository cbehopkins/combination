package combination

type ToDoChan chan []Position
type ToDoChanR <-chan []Position
type GenericWorker struct {
	todo ToDoChanR
}

type CombinationGenericInterface interface {
	CopyElement(src int, dst int)
}

func NewGenericWorker(todo ToDoChanR) *GenericWorker {
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
		dstArray.CopyElement(int(src), int(dst))
	}
	return nil
}

type GenericCombination struct {
	resultChan ToDoChanR
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
		gc.copyFunc(i, int(v))
	}
	return nil
}

func (gc GenericCombination) NextSkipN(n int) error {
	if gc.copyFunc == nil {
		return MissingCopyFuncError
	}
	var resArray []Position
	var ok bool
	for i := 0; i < n; i++ {
		resArray, ok = <-gc.resultChan
		if !ok {
			return ItterationCompleteError
		}
	}
	for i, v := range resArray {
		gc.copyFunc(i, int(v))
	}
	return nil
}
