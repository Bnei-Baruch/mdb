package utils

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func ReadCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	return csv.NewReader(bufio.NewReader(f)).ReadAll()
}

func ParseCSVHeader(header []string) (map[string]int, error) {
	if len(header) == 0 {
		return nil, errors.New("Empty header")
	}

	h := make(map[string]int, len(header))
	for i, x := range header {
		h[strings.ToLower(x)] = i
	}

	return h, nil
}
