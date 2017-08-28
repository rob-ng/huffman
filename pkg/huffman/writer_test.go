package huffman

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {

	// Invariants
	// Header:
	// 1. First line should consist of space-separated numbers which should
	//    sum to length of unique inputs.
	// 2. Second line should exactly match unique inputs.
	// Body:
	// 1. Should exist

	testInputs := []struct {
		src string
		cb  Codebook
	}{
		{"1112222334",
			NewCanonicalCB(map[byte]float64{
				49: 0.3,
				50: 0.4,
				51: 0.2,
				52: 0.1,
			})},
	}

	for _, ti := range testInputs {
		var out bytes.Buffer

		hw := NewWriter(&out, ti.cb)
		for _, c := range ti.src {
			hw.Write([]byte(string(c)))
		}
		hw.Flush()

		header1, _ := out.ReadString('\n')
		header1Values := strings.Split(strings.Trim(header1, "\n"), " ")
		header1Total := 0
		for _, n := range header1Values {
			val, _ := strconv.Atoi(n)
			header1Total += val
		}
		if header1Total != len(ti.cb) {
			t.Errorf("Sum of values in first line of header should equal number of entries in codebook")
		}

		header2, _ := out.ReadString('\n')
		header2 = strings.Trim(header2, "\n")
		if len(header2) != len(ti.cb) {
			t.Errorf("Second line of header should be exactly as long as codebook")
		}
		_, err := out.ReadByte()
		if err != nil {
			t.Errorf("Body should not be empty")
		}
	}
}
