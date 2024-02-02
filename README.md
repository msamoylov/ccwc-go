# CCWC (Coding Challenges Word Count)

![build](https://github.com/msamoylov/ccwc-go/actions/workflows/go.yml/badge.svg) [![codecov](https://codecov.io/gh/msamoylov/ccwc-go/graph/badge.svg?token=Q4C78M9A53)](https://codecov.io/gh/msamoylov/ccwc-go)

CCWC is a custom implementation of a word count utility for the challenge from [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-wc). 
This tool provides functionalities to count bytes in text files, with a focus on performance and accuracy.

## Features

- Count the number of bytes in each input file.
- Count the number of characters in each input file.
- Count the number of words in each input file.
- Count the number of lines in each input file.

## Build

To build CCWC, follow these steps:

1. Ensure you have [Go installed](https://golang.org/dl/) on your system.
2. Clone the repository:
   ```bash
   git clone https://github.com/msamoylov/ccwc-go.git
   ```
3. Navigate to the cloned directory:
   ```bash
   cd ccwc-go
   ```
4. Build the project:
   ```bash
   make build
   ```

## Usage

```bash
# Generate a file containing 100K lines:
yes ccwc | head -n 100000 > 100k.txt
```

After building, you can use CCWC as follows:

```bash
# Count lines, words and bytes in files:
./ccwc 100k.txt
```
