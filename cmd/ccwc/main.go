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
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/msamoylov/ccwc/internal/wc"
)

var (
	printBytes bool
	printChars bool
	printLines bool
	printWords bool
	printHelp  bool
)

func processFlags() {
	flag.BoolVar(&printBytes, "c", false, "print the byte counts")
	flag.BoolVar(&printChars, "m", false, "print the character counts")
	flag.BoolVar(&printLines, "l", false, "print the newline counts")
	flag.BoolVar(&printWords, "w", false, "print the word counts")
	flag.BoolVar(&printHelp, "help", false, "display this help and exit")

	flag.Parse()

	// By default, wc prints three counts: the newline, words and byte counts.
	if !printBytes && !printLines && !printWords && !printChars {
		printBytes = true
		printLines = true
		printWords = true
	}
}

func help() {
	const usage = `Usage: ccwc [OPTION]... [FILE]...
Print newline, word, and byte counts for each FILE, and a total line if
more than one FILE is specified.  A word is a non-zero-length sequence of
printable characters delimited by white space.

With no FILE, or when FILE is -, read standard input.

The options below may be used to select which counts are printed, always in
the following order: newline, word, character, byte.
  -c            print the byte counts
  -m            print the character counts
  -l            print the newline counts
  -w            print the word counts
      --help    display this help and exit
`
	fmt.Print(usage)
}

func main() {
	processFlags()

	if printHelp {
		help()
		return
	}

	cfg := &wc.Config{
		CountBytes: printBytes,
		CountChars: printChars,
		CountLines: printLines,
		CountWords: printWords,
	}

	processor := wc.NewProcessor(cfg)

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	defer func() {
		if ferr := writer.Flush(); ferr != nil {
			fmt.Fprintln(os.Stderr, ferr)
		}
	}()

	total := struct {
		bytes     int
		chars     int
		lines     int
		words     int
		processed int
	}{}

	files := flag.Args()
	if len(files) == 0 { // Process stdin
		files = append(files, "-")
	}

	for _, file := range files {
		var f *os.File
		var err error

		switch file {
		case "":
			fmt.Fprintln(os.Stderr, "file name cannot be zero-length")
			continue
		case "-":
			f = os.Stdin
		default:
			f, err = os.Open(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}

		if err = processor.Analyze(f); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		writeCounts(writer, processor.Lines(), processor.Words(), processor.Chars(), processor.Bytes(), f.Name())

		total.processed++
		total.lines += processor.Lines()
		total.words += processor.Words()
		total.chars += processor.Chars()
		total.bytes += processor.Bytes()
	}

	if total.processed > 1 {
		writeCounts(writer, total.lines, total.words, total.chars, total.bytes, "total")
	}
}

func writeCounts(w *tabwriter.Writer, lines, words, chars, bytes int, name string) {
	if w == nil {
		return
	}

	var out []string
	if printLines {
		out = append(out, strconv.Itoa(lines))
	}
	if printWords {
		out = append(out, strconv.Itoa(words))
	}
	if printChars {
		out = append(out, strconv.Itoa(chars))
	}
	if printBytes {
		out = append(out, strconv.Itoa(bytes))
	}

	line := strings.Join(out, "\t") + "\t"

	if name != "/dev/stdin" {
		line += "\t" + name
	}

	if _, err := fmt.Fprintln(w, line); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
