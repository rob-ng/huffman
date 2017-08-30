package huffman

import (
	"strings"
	"testing"
)

func TestHeaderString(t *testing.T) {
	h := &Header{
		cb: codebook{
			&cbEntry{unit: 'a', code: 0, codeLen: 1},
			&cbEntry{unit: 'b', code: 2, codeLen: 2},
			&cbEntry{unit: 'c', code: 3, codeLen: 2},
		},
		numUnits: 10,
	}

	h.String()
}

// Will want to also check that no code is prefix of any other
func TestNewHeader(t *testing.T) {
	testInputs := []struct {
		unitWeights map[byte]float64
		numUnits    int
	}{
		{
			map[byte]float64{
				1: 0.1,
				2: 0.4,
				3: 0.2,
				4: 0.3,
			},
			10,
		},
	}

	for _, ti := range testInputs {
		h := NewHeader(ti.unitWeights, ti.numUnits)
		if h.numUnits != ti.numUnits {
			t.Errorf("Num units should match input. Was: %d, Expected: %d\n",
				h.numUnits, ti.numUnits)
		}
		for i := 0; i < len(h.cb)-1; i++ {
			first := h.cb[i]
			second := h.cb[i+1]
			if first.code >= second.code {
				t.Errorf("Code should be strictly increasing. i: %d, i+1: %d\n",
					first.code, second.code)
			} else if first.codeLen > second.codeLen {
				t.Errorf("Code length should be strictly increasing. i: %d, i+1: %d\n",
					first.codeLen, second.codeLen)
			} else if first.codeLen != second.codeLen &&
				ti.unitWeights[first.unit] < ti.unitWeights[second.unit] {
				t.Errorf("Larger weight should be correlated with shorter code. i: %d, i+1: %d\n",
					ti.unitWeights[first.unit], ti.unitWeights[second.unit])
			}
		}
	}
}

func TestNewHeaderFromReader(t *testing.T) {
	testInputs := []struct {
		src string
	}{
		{src: "thisisatestsourcestring"},
	}
	for _, ti := range testInputs {
		r := strings.NewReader(ti.src)
		h, err := NewHeaderFromReader(r)
		if err != nil {
			t.Errorf("Should have succesfully created new Header. Error: %s\n", err.Error())
		}
		if h.numUnits != len(ti.src) {
			t.Errorf("Num units should match input. Was: %d, Expected: %d\n",
				h.numUnits, len(ti.src))
		}
		for i := 0; i < len(h.cb)-1; i++ {
			first := h.cb[i]
			second := h.cb[i+1]
			if first.code >= second.code {
				t.Errorf("Code should be strictly increasing. i: %d, i+1: %d\n",
					first.code, second.code)
			} else if first.codeLen > second.codeLen {
				t.Errorf("Code length should be strictly increasing. i: %d, i+1: %d\n",
					first.codeLen, second.codeLen)
			}
		}
	}
}

func TestDeriveHeader(t *testing.T) {
	testInputs := []struct {
		desc     string
		cbLen    int
		numUnits int
	}{
		{
			"1 1 2\r\nabcd\r\n10\r\n\r\n",
			4,
			10,
		},
		// Should work if '\n' is a unit
		{
			"1 1 2\r\na\ned\r\n10\r\n\r\n",
			4,
			10,
		},
	}

	for _, ti := range testInputs {
		h := DeriveHeader(ti.desc)
		if h.numUnits != ti.numUnits {
			t.Errorf("Num units should match input. Was: %d, Expected: %d\n",
				h.numUnits, ti.numUnits)
		}
		if len(h.cb) != ti.cbLen {
			t.Errorf("Incorrect codebook length. Was: %d, Expected: %d\n",
				len(h.cb), ti.cbLen)
		}
		for i := 0; i < len(h.cb)-1; i++ {
			first := h.cb[i]
			second := h.cb[i+1]
			if first.code >= second.code {
				t.Errorf("Code should be strictly increasing. i: %d, i+1: %d\n",
					first.code, second.code)
			} else if first.codeLen > second.codeLen {
				t.Errorf("Code length should be strictly increasing. i: %d, i+1: %d\n",
					first.codeLen, second.codeLen)
			}
		}
	}
}
