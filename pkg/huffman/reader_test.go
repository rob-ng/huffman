package huffman

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	utf8src := "©©©»»»»かかπ"
	br := bufio.NewReader(strings.NewReader(utf8src))
	utf8H, _ := NewHeaderFromReader(br)
	testInputs := []struct {
		src string
		//cb  Codebook
		h *Header
	}{
		{
			"1112222334",
			NewHeader(map[byte]float64{
				49: 0.3,
				50: 0.4,
				51: 0.2,
				52: 0.1,
			}, 10),
		}, {
			utf8src,
			utf8H,
		},
	}

	for _, ti := range testInputs {
		var out bytes.Buffer

		hw := NewWriter(&out, ti.h)
		hw.Write([]byte(ti.src))
		hw.Flush()

		r := bufio.NewReader(&out)
		hr := NewReader(r)

		out2 := make([]byte, ti.h.numUnits+10)
		hr.Read(out2)

		inBytes := []byte(ti.src)
		for i, _ := range inBytes {
			if out2[i] != inBytes[i] {
				t.Errorf("Decoded data should equal source data. Was: %d, Expected: %d\n",
					out2[i], ti.src[i])
			}
		}
		for i := ti.h.numUnits; i < ti.h.numUnits+10; i++ {
			if out2[i] != 0 {
				t.Errorf("Should only have read %d bytes to out", ti.h.numUnits)
			}
		}

		//fmt.Printf("%s\n", out2)
	}
}
