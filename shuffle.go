package combination

import "log"

// We want to genericise something we do in countdown package
// That is here suppose we have [0, 1, 2, 3]
// We want:
// [[0], [1, 2, 3]]
// [[1], [0, 2, 3]]
// [[2], [0, 1, 3]]
// [[3], [0, 1, 2]]
// [[0, 1], [2, 3]]
// [[0, 2], [1, 3]]
// [[0, 3], [1, 2]]
// That should give all combinations as for our purposes:
// [[0, 1], [2, 3]] === [[2, 3], [0, 1]] === [[1, 0], [2, 3]] === [[1, 0], [3, 2]]
// See the test Examples for what we are tryng to generate

func oneLevel(sizeUp, stopAt int, sc chan<- ToDo, prevArray ToDo) {
	genSize := len(prevArray) + 1
	if stopAt < genSize {
		// Spawn nothing longer than stopAt
		return
	}
	for i := sizeUp; i >= 0; i-- {
		tmpAr := make(ToDo, genSize)
		copy(tmpAr, prevArray)
		tmpAr[len(prevArray)] = Position(i)

		// Send the array of length N
		// e.g. [4 3]
		sc <- tmpAr
		nextSa := i - 1
		if nextSa < stopAt {
			// Decide where to stop at:
			// note for e.g.  len 5, [4 2]
			// we don't want to then go on and generate [4 2 1]
			// as that is covered by [3 0]
			// or if len 6 stop at [4 3 2]
			// as [4 3 2 1] is covered by [5 0]
			stopAt--
		}

		// Generate and send Length N arrays basedoff this
		// particular array. e.g. from tmpAr = [4 3]
		// [4 3 2]
		// [4 3 1]
		// [4 3 0]
		oneLevel(nextSa, stopAt, sc, tmpAr)
	}
}

func newShuffleChan(sizeUp int) <-chan ToDo {
	sc := make(chan ToDo)
	go func() {
		oneLevel(sizeUp-1, sizeUp, sc, ToDo{})
		close(sc)
	}()
	return sc
}

type ShuffleGen struct {
	size int
}

func (td ToDo) contains(val Position) bool {
	for _, i := range td {
		if i == val {
			return true
		}
	}
	return false
}
func (sg ShuffleGen) remains(in ToDo) ToDo {

	fillLen := sg.size - len(in)
	destArray := make([]Position, 0, fillLen)
	for i := 0; i < sg.size; i++ {
		if !in.contains(Position(i)) {
			destArray = append(destArray, Position(i))
		}
	}
	if len(destArray) != fillLen {
		log.Println("Error with fill length not matching filled array")
		return nil
	}
	return destArray
}
func (sg ShuffleGen) splitter(inputChan <-chan ToDo) chan []ToDo {
	resultChan := make(chan []ToDo)
	go func() {
		for ip := range inputChan {
			remains := sg.remains(ip)
			resultChan <- []ToDo{ip, remains}
		}
		close(resultChan)
	}()
	return resultChan
}
func NewShuffleGen(size int) <-chan []ToDo {
	sg := ShuffleGen{size}
	sc := newShuffleChan(size)
	return sg.splitter(sc)
}
