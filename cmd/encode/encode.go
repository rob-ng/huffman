package main

import (
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
	r := io.TeeReader(inFile, &buf)
	header, err := huffman.NewHeaderFromReader(r)
	if err != nil {
		panic(err)
	}
	hw := huffman.NewWriter(outFile, header)
	hw.Write(buf.Bytes())
	hw.Flush()
}
