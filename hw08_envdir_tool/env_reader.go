package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	environment := Environment{}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("fail to get info: %w", err)
		}
		if entry.IsDir() {
			return nil, fmt.Errorf("incorrect file path: %w", err)
		}
		if strings.ContainsAny(entry.Name(), "=") {
			return nil, fmt.Errorf("incorrect file name: %w", err)
		}

		file, err := os.Open(path.Join(dir, entry.Name()))
		if err != nil {
			fmt.Println("Error opening file:", err)
			return nil, err
		}
		defer file.Close()

		firstLine, err := readFirstLine(file)
		if err != nil {
			return nil, err
		}

		firstLine = normalize(firstLine)
		env := EnvValue{Value: firstLine, NeedRemove: info.Size() == 0}

		environment[entry.Name()] = env
	}

	return environment, nil
}

func readFirstLine(reader io.Reader) (string, error) {
	var firstLine string

	scanner := bufio.NewScanner(reader)

	if scanner.Scan() {
		firstLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return firstLine, nil
}

func normalize(str string) string {
	str = strings.TrimRight(str, " \t")

	if str == " " {
		str = ""
	}

	nullIndex := strings.Index(str, "\x00")

	if nullIndex > -1 {
		str = str[:nullIndex]
	}

	return str
}
