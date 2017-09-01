package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rob-ng/huffman/pkg/huffman"
)

var decompressFlag = flag.Bool("decompress", false, "decompress file")

// Set short versions of flags
func init() {
	flag.BoolVar(decompressFlag, "d", false, "decompress file")
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalln(err)
		}
	}()

	flag.Parse()
	args := flag.Args()

	source := os.Stdin
	output := os.Stdout

	if len(args) >= 1 && strings.Compare(args[0], "-") != 0 {
		if source, err = os.Open(args[0]); err != nil {
			return
		}
		defer func() {
			if err = source.Close(); err != nil {
				return
			}
		}()
	}

	if len(args) >= 2 && strings.Compare(args[1], "-") != 0 {
		if output, err = os.Create(args[1]); err != nil {
			return
		}
		defer func() {
			if err = output.Close(); err != nil {
				return
			}
		}()
	}

	if *decompressFlag {
		r := huffman.NewReader(bufio.NewReader(source))
		w := output
		buf := make([]byte, 1024)
		shouldExit := false
		var n int
		for {
			if n, err = r.Read(buf); err != nil {
				if err == io.EOF {
					err = nil
					shouldExit = true
				} else {
					return
				}
			}

			if n == 0 {
				break
			}

			if _, err = w.Write(buf[:n]); err != nil {
				return
			}

			if shouldExit {
				break
			}
		}
	} else {
		var buf bytes.Buffer
		tr := io.TeeReader(bufio.NewReader(source), &buf)
		var header *huffman.Header
		if header, err = huffman.NewHeaderFromReader(tr); err != nil {
			return
		}
		w := huffman.NewWriter(output, header)
		w.Write(buf.Bytes())
		w.Flush()
	}
}
