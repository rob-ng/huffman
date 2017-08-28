package huffman

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	utf8src := `©©©»»»»かかπ`
	br := bufio.NewReader(strings.NewReader(utf8src))
	utf8WM, _ := MakeWeightMap(br)
	testInputs := []struct {
		src string
		cb  Codebook
	}{
		{
			"1112222334",
			NewCanonicalCB(map[byte]float64{
				49: 0.3,
				50: 0.4,
				51: 0.2,
				52: 0.1,
			}),
		}, {
			utf8src,
			NewCanonicalCB(utf8WM),
		},
	}

	for _, ti := range testInputs {
		var out bytes.Buffer

		hw := NewWriter(&out, ti.cb)
		for _, c := range ti.src {
			hw.Write([]byte(string(c)))
		}
		hw.Flush()

		r := bufio.NewReader(&out)
		hr := NewReader(r)

		out2 := make([]byte, 10)
		hr.Read(out2)

		for i, _ := range out2 {
			if out2[i] != ti.src[i] {
				t.Errorf("Decoded data should equal source data")
			}
		}
	}
}
