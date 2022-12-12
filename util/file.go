package util

import (
	"bufio"
	"io"
	"os"
)

func ReadLine2Slice(file string) ([]string, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var lines []string
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}
	return lines, nil

}
