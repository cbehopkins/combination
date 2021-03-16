package combination

import "errors"

var ItterationCompleteError = errors.New("Exhausted itterator")
var MissingCopyFuncError = errors.New("No copy func supplied")

type Combination struct {
	len int
}

func NewCombination(length int) *Combination {
	v := new(Combination)
	v.len = length
	return v
}
func (c Combination) ToChannel(resultChan chan<- []int) {
	referenceArray := c.makeReferenceArray(c.len)
	c.toChannelInput(referenceArray, resultChan, c.len)
	close(resultChan)
}
func (c Combination) makeReferenceArray(length int) []int {
	referenceArray := make([]int, length)
	for i := range referenceArray {
		referenceArray[i] = i
	}
	return referenceArray
}
func (c Combination) toChannelInput(referenceArray []int, resultChan chan<- []int, srcLen int) {

	for i := 0; i < srcLen; i++ {
		tmpArray := make([]int, 0, srcLen-1)
		tmpArray = append(tmpArray, referenceArray[:i]...)
		tmpArray = append(tmpArray, referenceArray[i+1:]...)
		resultChan <- tmpArray
	}
}

func NewCombChannel(length int) <-chan []int {
	cmb := NewCombination(length)
	resultChan := make(chan []int)
	go cmb.ToChannel(resultChan)
	return resultChan
}

// NewCombChannelLen goes from a source array length
// to a destination length
func NewCombChannelLen(src, dst int) <-chan []int {
	if src <= dst {
		return nil
	}
	c := NewCombination(src)
	ra := c.makeReferenceArray(src)

	srcChan := make([]chan []int, 1)
	srcChan[0] = make(chan []int)

	go func() {
		srcChan[0] <- ra
		close(srcChan[0])
	}()
	currentLen := src
	i := 0
	for currentLen > dst {
		currentLen--
		scc := make(chan []int)
		srcChan = append(srcChan, scc)
		go c.bob(srcChan[i], srcChan[i+1])
		i++
	}
	return srcChan[i]
}
func (c Combination) bob(srcChan, dstChan chan []int) {
	for ra := range srcChan {
		c.toChannelInput(ra, dstChan, len(ra))
	}
	close(dstChan)

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
