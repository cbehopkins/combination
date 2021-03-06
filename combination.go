package combination

import (
	"errors"
	"fmt"
)

var ErrItterationComplete = errors.New("exhausted itterator")
var ErrMissingCopyFunc = errors.New("no copy func supplied")

type Position int
type ToDo []Position

var ErrToDoTooLong = errors.New("ToDo is too long for hash")
var ErrToDoOutOfRange = errors.New("To Do is out of range")

func (td ToDo) Hash() (uint64, error) {
	if len(td) > 8 {
		return 0, ErrToDoTooLong
	}
	var sum uint64
	for i := 0; i < len(td); i++ {
		if td[i] > 255 {
			return 0, fmt.Errorf("Got Position %d which is improbable, and not Hashable %w", td[i], ErrToDoOutOfRange)
		}
		sum <<= 8
		sum += uint64(td[i] & 255)
	}

	return sum, nil
}

type Combination struct {
	len int
}

func NewCombination(length int) *Combination {
	v := new(Combination)
	v.len = length
	return v
}
func (c Combination) ToChannel(resultChan chan<- ToDo) {
	referenceArray := c.makeReferenceArray(c.len)
	c.toChannelInput(referenceArray, resultChan, c.len)
	close(resultChan)
}
func (c Combination) makeReferenceArray(length int) []Position {
	referenceArray := make(ToDo, length)
	for i := range referenceArray {
		referenceArray[i] = Position(i)
	}
	return referenceArray
}
func (c Combination) toChannelInput(referenceArray ToDo, resultChan chan<- ToDo, srcLen int) {

	for i := 0; i < srcLen; i++ {
		tmpArray := make(ToDo, 0, srcLen-1)
		tmpArray = append(tmpArray, referenceArray[:i]...)
		tmpArray = append(tmpArray, referenceArray[i+1:]...)
		resultChan <- tmpArray
	}
}

func NewCombChannel(length int) <-chan ToDo {
	cmb := NewCombination(length)
	resultChan := make(chan ToDo, 16)
	go cmb.ToChannel(resultChan)
	return resultChan
}

// NewCombChannelLen goes from a source array length
// to a destination length
func NewCombChannelLen(src, dst int) <-chan ToDo {
	if src <= dst {
		return nil
	}
	c := NewCombination(src)
	ra := c.makeReferenceArray(src)

	srcChan := make([]chan ToDo, 1)
	srcChan[0] = make(chan ToDo)

	go func() {
		srcChan[0] <- ra
		close(srcChan[0])
	}()
	currentLen := src
	i := 0
	for currentLen > dst {
		currentLen--
		var scc chan ToDo
		if currentLen == dst {
			scc = make(chan ToDo, 16)
		} else {
			scc = make(chan ToDo)
		}
		srcChan = append(srcChan, scc)
		go c.copyWorker(srcChan[i], srcChan[i+1])
		i++
	}
	return srcChan[i]
}
func (c Combination) copyWorker(srcChan, dstChan chan ToDo) {
	for ra := range srcChan {
		c.toChannelInput(ra, dstChan, len(ra))
	}
	close(dstChan)
}
