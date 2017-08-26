package huffman

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestNewCanonincalCB(t *testing.T) {
	testInputs := make([]map[byte]float64, 10)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i := 0; i < 10; i++ {
		total := 0.0
		testInputs[i] = make(map[byte]float64)
		for total < 1.0 {
			weight := r1.Float64()
			bite := byte(r1.Intn(256))
			total += weight
			testInputs[i][bite] = weight
		}
	}

	for _, ti := range testInputs {
		cb := NewCanonicalCB(ti)
		// 3 invariants
		// 1. Value of code is strictly increasing
		// 2. Code length is increasing
		// 3. Between different code lengths, longer lengths should
		//    be associated with smaller weights
		for i := 0; i < len(cb)-1; i++ {
			first, _ := strconv.ParseInt(cb[i].code, 2, 8)
			second, _ := strconv.ParseInt(cb[i+1].code, 2, 8)
			if first >= second {
				t.Errorf("Code should be increasing")
			} else if cb[i].codeLen > cb[i+1].codeLen {
				t.Errorf("Code len should be increasing")
			} else if cb[i].codeLen != cb[i+1].codeLen &&
				ti[cb[i].unit] < ti[cb[i+1].unit] {
				t.Errorf("Frequency decrasing!")
			}
		}
	}
}
