package main

import (
	"bufio"
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

	br := bufio.NewReader(inFile)
	hr := huffman.NewReader(br)

	buf := make([]byte, 1024)
	shouldExit := false
	for {
		n, err := hr.Read(buf)
		if err != nil {
			if err == io.EOF {
				shouldExit = true
			} else {
				panic(err)
			}
		}
		if n == 0 {
			break
		}

		if _, err := outFile.Write(buf[:n]); err != nil {
			panic(err)
		}

		if shouldExit {
			break
		}
	}
}
