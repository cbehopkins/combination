package combination

import "errors"

var ItterationCompleteError = errors.New("Exhausted itterator")
var MissingCopyFuncError = errors.New("No copy func supplied")

type Position int
type PosChan chan []Position

type Combination struct {
	len int
}

func NewCombination(length int) *Combination {
	v := new(Combination)
	v.len = length
	return v
}
func (c Combination) ToChannel(resultChan chan<- []Position) {
	referenceArray := c.makeReferenceArray(c.len)
	c.toChannelInput(referenceArray, resultChan, c.len)
	close(resultChan)
}
func (c Combination) makeReferenceArray(length int) []Position {
	referenceArray := make([]Position, length)
	for i := range referenceArray {
		referenceArray[i] = Position(i)
	}
	return referenceArray
}
func (c Combination) toChannelInput(referenceArray []Position, resultChan chan<- []Position, srcLen int) {

	for i := 0; i < srcLen; i++ {
		tmpArray := make([]Position, 0, srcLen-1)
		tmpArray = append(tmpArray, referenceArray[:i]...)
		tmpArray = append(tmpArray, referenceArray[i+1:]...)
		resultChan <- tmpArray
	}
}

func NewCombChannel(length int) <-chan []Position {
	cmb := NewCombination(length)
	resultChan := make(chan []Position)
	go cmb.ToChannel(resultChan)
	return resultChan
}

// NewCombChannelLen goes from a source array length
// to a destination length
func NewCombChannelLen(src, dst int) <-chan []Position {
	if src <= dst {
		return nil
	}
	c := NewCombination(src)
	ra := c.makeReferenceArray(src)

	srcChan := make([]chan []Position, 1)
	srcChan[0] = make(chan []Position)

	go func() {
		srcChan[0] <- ra
		close(srcChan[0])
	}()
	currentLen := src
	i := 0
	for currentLen > dst {
		currentLen--
		scc := make(chan []Position)
		srcChan = append(srcChan, scc)
		go c.bob(srcChan[i], srcChan[i+1])
		i++
	}
	return srcChan[i]
}
func (c Combination) bob(srcChan, dstChan chan []Position) {
	for ra := range srcChan {
		c.toChannelInput(ra, dstChan, len(ra))
	}
	close(dstChan)

}
