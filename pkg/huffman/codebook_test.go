package huffman

import (
	"math/rand"
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
		// 4 invariants
		// 1. Value of code is strictly increasing
		// 2. Code length is increasing
		// 3. Between different code lengths, longer lengths should
		//    be associated with smaller weights
		// 4. No code is prefix of any other code
		for i := 0; i < len(cb)-1; i++ {
			first := cb[i].code
			second := cb[i+1].code
			if first >= second {
				t.Errorf("Code should be strictly increasing")
			} else if cb[i].codeLen > cb[i+1].codeLen {
				t.Errorf("Code len should be increasing")
			} else if cb[i].codeLen != cb[i+1].codeLen &&
				ti[cb[i].unit] < ti[cb[i+1].unit] {
				t.Errorf("Code length should be correlated with weight")
			}
		}
	}
}

func TestDeriveCanonicalCB(t *testing.T) {
	deriveTests := []struct {
		header   *Header
		codes    []int
		codeLens []int
	}{
		{
			&Header{[]int{1, 1, 2}, []byte("abcd")},
			[]int{0, 2, 6, 7},
			[]int{1, 2, 3, 3},
		},
	}

	for _, dt := range deriveTests {
		cb := DeriveCanonicalCB(dt.header)
		for i, e := range cb {
			if dt.codes[i] != e.code {
				t.Errorf("Incorrect code! Was %d, expected %d", e.code, dt.codes[i])
			}
			if dt.codeLens[i] != e.codeLen {
				t.Errorf("Incorrect code length! Was %d, expected %d", e.codeLen, dt.codeLens[i])
			}
		}
	}
}
