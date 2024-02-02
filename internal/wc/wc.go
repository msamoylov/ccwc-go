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

// Package wc implements the wc tool from GNU Coreutils.
// See https://www.gnu.org/software/coreutils/manual/html_node/wc-invocation.html
package wc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Config provides configuration for Processor.
type Config struct {
	CountLines bool
	CountWords bool
	CountChars bool
	CountBytes bool
}

// Processor handles the counting of bytes, characters, lines, and words in a text stream.
type Processor struct {
	lines int // Count of lines
	words int // Count of words
	chars int // Count of characters
	bytes int // Count of bytes

	cfg *Config
}

// NewProcessor creates a new Processor.
func NewProcessor(cfg *Config) *Processor {
	if cfg == nil {
		cfg = &Config{}
	}
	return &Processor{cfg: cfg}
}

func (p *Processor) reset() {
	p.bytes = 0
	p.chars = 0
	p.lines = 0
	p.words = 0
}

// Bytes returns count of bytes.
func (p *Processor) Bytes() int {
	return p.bytes
}

// Chars returns count of chars.
func (p *Processor) Chars() int {
	return p.chars
}

// Lines returns count of lines.
func (p *Processor) Lines() int {
	return p.lines
}

// Words returns count of words.
func (p *Processor) Words() int {
	return p.words
}

// Analyze counts the number of bytes, characters, words, and newlines in the file or standard input.
func (p *Processor) Analyze(f *os.File) error {
	if f == nil {
		return fmt.Errorf("input cannot be nil")
	}

	var err error
	defer func() {
		if closeErr := f.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	p.reset() // Always reset just in case

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("read %v: is a directory", f.Name())
	}

	// When counting only bytes in a file, save some line- and word-counting overhead.
	if f.Name() != "/dev/stdin" && p.cfg.CountBytes && !p.cfg.CountChars && !p.cfg.CountLines && !p.cfg.CountWords {
		p.bytes = int(fi.Size())
		return nil
	}

	lineReader := bufio.NewReader(f)
	for {
		line, err := lineReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		p.lines++
		p.bytes += len(line)

		if p.cfg.CountChars {
			for range line {
				p.chars++
			}
		}

		if p.cfg.CountWords {
			wordScanner := bufio.NewScanner(strings.NewReader(line))
			wordScanner.Split(bufio.ScanWords)
			for wordScanner.Scan() {
				p.words++
			}
		}
	}

	return nil
}
