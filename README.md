# ccwc
Build Your Own `wc` Tool.

# CCWC (Coding Challenges Word Count)

## Introduction

CCWC is a custom implementation of a word count utility for the challenge from [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-wc). 
This tool provides functionalities to count bytes in text files, with a focus on performance and accuracy.

## Features

- Count the number of bytes in each input file.

## Installation

To install CCWC, follow these steps:

1. Ensure you have [Go installed](https://golang.org/dl/) on your system.
2. Clone the repository:
   ```bash
   git clone https://github.com/msamoylov/ccwc.git
   ```
3. Navigate to the cloned directory:
   ```bash
   cd ccwc
   ```
4. Build the project:
   ```bash
   make build
   ```

## Usage

After installing, you can use CCWC as follows.

To count bytes in files:

```bash
./ccwc -c file1.txt file2.txt
```