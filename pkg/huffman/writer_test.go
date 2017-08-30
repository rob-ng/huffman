package huffman

import (
	"bytes"
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
		h   *Header
	}{
		{
			"1112222334",
			NewHeader(map[byte]float64{
				49: 0.3,
				50: 0.4,
				51: 0.2,
				52: 0.1,
			}, 10),
		},
	}

	for _, ti := range testInputs {
		var out bytes.Buffer
		hw := NewWriter(&out, ti.h)
		hw.Write([]byte(ti.src))
		hw.Flush()

		if out.Len() == 0 {
			t.Errorf("Encoded data should have been written to buffer")
		}

		//fmt.Printf("%v\n", out.Bytes())

		/*header1, _ := out.ReadString('\n')
		header1Values := strings.Split(strings.Trim(header1, "\n"), " ")
		header1Total := 0
		for _, n := range header1Values {
			val, _ := strconv.Atoi(n)
			header1Total += val
		}
		if header1Total != len(ti.h.cb) {
			t.Errorf("Sum of values in first line of header should equal number of entries in codebook. Was: %d, Expected; %d",
				header1Total, len(ti.h.cb))
		}

		header2, _ := out.ReadString('\n')
		header2 = strings.Trim(header2, "\n")
		if len(header2) != len(ti.h.cb) {
			t.Errorf("Second line of header should be exactly as long as codebook")
		}
		header3, _ := out.ReadString('\n')
		header3 = strings.Trim(header3, "\n")
		if len(header3) == 0 {
			t.Errorf("Third line should contain unit count")
		}

		_, err := out.ReadByte()
		if err != nil {
			t.Errorf("Body should not be empty")
		}*/
	}
}
