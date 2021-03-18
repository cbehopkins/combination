package combination

import (
	"fmt"
	"testing"
)

type resultChecker map[string]int

func (rc resultChecker) add(ip []int) {
	key := fmt.Sprint(ip)
	v := rc[key]
	v++
	rc[key] = v
}
func (rc resultChecker) tst(ip []Position) bool {
	key := fmt.Sprint(ip)
	v, ok := rc[key]
	v--
	if v <= 0 {
		delete(rc, key)
		v = 0
	} else {

		rc[key] = v
	}
	return ok
}
func runSimpleTest(length int, t *testing.T, expectedResults resultChecker) (missingCnt int) {
	tmp := NewCombination(length)
	resultChan := make(chan ToDo)
	go tmp.ToChannel(resultChan)
	for v := range resultChan {
		t.Log("Got:", v)
		if !expectedResults.tst(v) {
			missingCnt++
			t.Log("The value", v, "is not in our expected list")
		}
	}
	return
}
func TestBasicCombination(t *testing.T) {
	// As an example:
	// If my input width is 3, and output width is 2
	// I'd expect to get
	// [0, 1]
	// [1, 2]
	// [0, 2]
	// In some order
	expectedResults := make(resultChecker)
	expectedResults.add([]int{0, 1})
	expectedResults.add([]int{1, 2})
	expectedResults.add([]int{0, 2})
	missingCnt := runSimpleTest(3, t, expectedResults)
	if len(expectedResults) != 0 {
		t.Error("We didn't get all expected values")
	}
	if missingCnt != 0 {
		t.Error("There should have been no missing values")
	}
}

func TestExtraExpectedResults(t *testing.T) {
	// Exactly as above, but make sure we fail
	// If we are missing expected results
	expectedResults := make(resultChecker)
	expectedResults.add([]int{0, 1})
	expectedResults.add([]int{1, 2})
	expectedResults.add([]int{0, 2})
	expectedResults.add([]int{0, 3})

	missingCnt := runSimpleTest(3, t, expectedResults)
	if len(expectedResults) != 1 {
		t.Error("We didn't get all expected values")
	}
	if missingCnt != 0 {
		t.Error("There should have been no missing values")
	}
}

func TestMissingExpectedResults(t *testing.T) {
	// Exactly as above, but make sure we fail
	// If we are missing expected results
	expectedResults := make(resultChecker)
	expectedResults.add([]int{0, 1})
	expectedResults.add([]int{0, 2})
	length := 3
	missingCnt := runSimpleTest(length, t, expectedResults)
	if len(expectedResults) != 0 {
		t.Error("We didn't get all expected values")
	}
	if missingCnt != 1 {
		t.Error("There should have been one missing value")
	}
}
func runAutoTest(t *testing.T, expectedResults resultChecker, src, dst int) (missingCnt int) {
	resultChan := NewCombChannelLen(src, dst)
	for v := range resultChan {
		t.Log("Got:", v)
		if !expectedResults.tst(v) {
			missingCnt++
			t.Log("The value", v, "is not in our expected list")
		}
	}
	return
}
func TestAutoCombination_3_2(t *testing.T) {
	// As TestBasicCombination, but use the auto function
	expectedResults := make(resultChecker)
	expectedResults.add([]int{0, 1})
	expectedResults.add([]int{1, 2})
	expectedResults.add([]int{0, 2})
	missingCnt := runAutoTest(t, expectedResults, 3, 2)
	if len(expectedResults) != 0 {
		t.Error("We didn't get all expected values")
	}
	if missingCnt != 0 {
		t.Error("There should have been no missing values")
	}
}
func TestAutoCombination_4_3(t *testing.T) {
	// As TestAutoCombination_3_2, but larger width
	expectedResults := make(resultChecker)
	expectedResults.add([]int{1, 2, 3})
	expectedResults.add([]int{0, 2, 3})
	expectedResults.add([]int{0, 1, 3})
	expectedResults.add([]int{0, 1, 2})
	missingCnt := runAutoTest(t, expectedResults, 4, 3)
	if len(expectedResults) != 0 {
		t.Error("We didn't get all expected values")
	}
	if missingCnt != 0 {
		t.Error("There should have been no missing values")
	}
}
func TestAutoCombination_4_2(t *testing.T) {
	// As TestAutoCombination_4_3, but now
	// a case that will show the recursion
	expectedResults := make(resultChecker)
	// Note each one of these triplets is a, expansion of Start
	expectedResults.add([]int{2, 3})
	expectedResults.add([]int{1, 3})
	expectedResults.add([]int{1, 2})

	expectedResults.add([]int{2, 3})
	expectedResults.add([]int{0, 3})
	expectedResults.add([]int{0, 2})

	expectedResults.add([]int{1, 3})
	expectedResults.add([]int{0, 3})
	expectedResults.add([]int{0, 1})

	expectedResults.add([]int{1, 2})
	expectedResults.add([]int{0, 2})
	expectedResults.add([]int{0, 1})

	missingCnt := runAutoTest(t, expectedResults, 4, 2)
	if len(expectedResults) != 0 {
		t.Error("We didn't get all expected values")
	}
	if missingCnt != 0 {
		t.Error("There should have been no missing values")
	}
}
