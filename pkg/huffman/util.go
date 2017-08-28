package huffman

import (
	"bufio"
	"io"
)

func MakeWeightMap(in *bufio.Reader) (weightMap map[byte]float64, err error) {
	weightMap = make(map[byte]float64)
	var n float64
	for {
		b, err := in.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		if _, ok := weightMap[b]; !ok {
			weightMap[b] = 0
		}
		weightMap[b]++
		n++
	}
	for i, _ := range weightMap {
		weightMap[i] = weightMap[i] / n
	}
	return weightMap, nil
}
