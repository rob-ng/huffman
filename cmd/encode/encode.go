package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/rob-ng/huffman/pkg/huffman"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Printf("Usage: %s inputFile outputFile\n", os.Args[0])
		os.Exit(1)
	}

	inFile, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	outFile, err := os.Create(args[1])
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	var buf bytes.Buffer
	br := bufio.NewReader(io.TeeReader(inFile, &buf))

	unitWeights, numUnits, _ := huffman.ProcessData(br)
	header := huffman.NewHeader(unitWeights, numUnits)
	hw := huffman.NewWriter(outFile, header)
	hw.Write(buf.Bytes())
	hw.Flush()
}
