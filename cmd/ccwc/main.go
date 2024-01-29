// Copyright (c) 2024 Michael Samoylov.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/msamoylov/ccwc/internal/file"
)

func main() {
	countBytes := flag.Bool("c", true, "Count the number of bytes in each input file.")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "No input specified.")
		os.Exit(1)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', tabwriter.AlignRight)

	var total int64
	for _, path := range files {
		if *countBytes {
			size, err := file.Size(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ccwc: %v\n", err)
				continue
			}

			output(writer, fmt.Sprintf("%d\t%s", size, path))
			total += size
		}
	}

	output(writer, fmt.Sprintf("%d\ttotal", total))

	if err := writer.Flush(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func output(w *tabwriter.Writer, line string) {
	if _, err := fmt.Fprintln(w, line); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
