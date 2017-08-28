package huffman

import (
	"bufio"
	"strings"
	"testing"
)

func TestMakeWeightMap(t *testing.T) {
	testInputs := []struct {
		src string
	}{
		{
			`©©©©»»»かかπ`,
		},
	}

	for _, ti := range testInputs {
		// Weights should sum to one
		br := bufio.NewReader(strings.NewReader(ti.src))
		wm, err := MakeWeightMap(br)
		if err != nil {
			t.Errorf("Should not have encountered error! Was: %s", err.Error())
		}
		sum := 0.0
		for _, v := range wm {
			sum += v
		}
		if sum != 1.0 {
			t.Errorf("Weights should sum to 1. Was: %d\n", sum)
		}
	}
}
