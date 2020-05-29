package main

import (
	"bufio"
	"io"
	"math/rand"
	"os"
	"strings"
)

func checkFileLines(path string) int {
	var i = 0

	file, err := os.Open(path)
	checkError(err)

	reader := bufio.NewReader(file)
	defer file.Close()

	for {
		_, err := reader.ReadString('\n')
		i++
		if err != nil && err == io.EOF {
			break
		}
		checkError(err)
	}
	return i
}

func parseFileToLines(path string) ([]string, int) {
	var lineNum = checkFileLines(path)
	var lines = make([]string, lineNum)
	var i = 0

	file, err := os.Open(path)
	checkError(err)

	reader := bufio.NewReader(file)
	defer file.Close()

	for {
		line, err := reader.ReadString('\n')
		lines[i] = strings.ReplaceAll(line, "\n", "")
		if err != nil && err == io.EOF {
			break
		}
		checkError(err)
		i++
	}
	return lines, lineNum
}

func phrasePicker(path string) string {
	var phrase string

	lines, lineNum := parseFileToLines(path)
	if lineNum <= 1 {
		phrase = lines[0]
	} else {
		phrase = lines[rand.Intn(lineNum-1)]
	}

	return phrase
}
