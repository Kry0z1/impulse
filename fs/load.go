package fs

import (
	"bufio"
	"fmt"
	"os"
)

func MustLoadInputFile(path string) *bufio.Reader {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed to open input file: %w", err))
	}

	return bufio.NewReader(f)
}

func OpenOutputFile(path string) (*bufio.Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %w", err)
	}

	return bufio.NewWriter(f), nil
}

func OpenLogFile(path string) (*bufio.Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return bufio.NewWriter(f), nil
}
