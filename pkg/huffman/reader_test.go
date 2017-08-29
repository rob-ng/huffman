package huffman

import (
	"bufio"
	"bytes"
	"testing"
)

func TestRead(t *testing.T) {
	//utf8src := "©©©»»»»かかπ"
	//br := bufio.NewReader(strings.NewReader(utf8src))
	//utf8WM, _ := MakeWeightMap(br)
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
		}, /*{
			utf8src,
			NewHeader(utf8WM, 10),
		},*/
	}

	for _, ti := range testInputs {
		var out bytes.Buffer

		hw := NewWriter(&out, ti.h)
		hw.Write([]byte(ti.src))
		hw.Flush()

		r := bufio.NewReader(&out)
		hr := NewReader(r)

		out2 := make([]byte, 10)
		hr.Read(out2)

		for i, _ := range out2 {
			if out2[i] != ti.src[i] {
				t.Errorf("Decoded data should equal source data. Was: %d, Expected; %d\n",
					out2[i], ti.src[i])
			}
		}
	}
}
