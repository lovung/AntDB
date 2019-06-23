package golib

import (
	"io/ioutil"
	"strconv"
)

// ReadWholeFile read whole file which have integers in every lines
func ReadWholeFile(filename string, arr *[]int) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var lasti, val int
	for i, v := range data {
		if v == '\n' {
			s := string(data[lasti:i])
			val, _ = strconv.Atoi(s)
			*arr = append(*arr, val)
			lasti = i + 1
		}
	}
	return nil
}

// WriteWholeFile write all data to a file once
func WriteWholeFile(filename string, arr []int) error {
	var buffer []byte
	for _, val := range arr {
		buffer = append(buffer, []byte(strconv.Itoa(val))...)
		buffer = append(buffer, byte('\n'))
	}

	return ioutil.WriteFile(filename, buffer, 0644)
}
