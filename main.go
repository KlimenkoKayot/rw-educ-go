package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func WriteString(s string, w io.Writer) error {
	_, err := w.Write([]byte(s))
	return err
}

func ReadString(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type UpperWriter struct {
	UpperString string
}

func (uw *UpperWriter) Write(p []byte) (int, error) {
	str := strings.ToUpper(string(p))
	uw.UpperString = str
	return len(p), nil
}

func Copy(r io.Reader, w io.Writer, n uint) error {
	buffer := make([]byte, n)
	_, err := r.Read(buffer)
	if err != nil {
		return err
	}
	for buffer[len(buffer)-1] == 0 {
		buffer = buffer[:len(buffer)-1]
	}
	_, err = w.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

func Contains(r io.Reader, seq []byte) (bool, error) {
	if len(seq) == 0 {
		return false, fmt.Errorf("UB, seq is empty")
	}
	buf := make([]byte, 1024*512)
	searchBuffer := []byte{}
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if n > 0 {
			searchBuffer = append(searchBuffer, buf...)
			if len(searchBuffer) >= len(seq) {
				if bytes.Contains(searchBuffer, seq) {
					return true, nil
				}
			}
			if len(searchBuffer) >= 2*len(seq) {
				searchBuffer = searchBuffer[len(searchBuffer)-len(seq):]
			}
		}
	}
	return false, nil
}

func CopyFilePart(inputFilename, outFileName string, startpos int) error {
	file, err := os.Open(inputFilename)
	if err != nil {
		return fmt.Errorf("open file error")
	}
	defer file.Close()

	out, err := os.Create(outFileName)
	if err != nil {
		return fmt.Errorf("create out file error")
	}
	file.Seek(int64(startpos), 0)

	buf := make([]byte, 128*1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read out file error")
		}
		out.Write(buf[:n])
	}
	return nil
}

func ModifyFile(filename string, pos int, val string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open file error")
	}
	defer file.Close()

	file.Seek(int64(pos), 0)
	file.WriteString(val)
	return nil
}

func ExtractLog(inputFileName string, start, end time.Time) ([]string, error) {
	file, err := os.Open(inputFileName)
	if err != nil {
		return nil, fmt.Errorf("open file error")
	}
	defer file.Close()

	result := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()

		year, _ := strconv.Atoi(str[6:10])
		month, _ := strconv.Atoi(str[3:5])
		day, _ := strconv.Atoi(str[0:2])
		// result = append(result, fmt.Sprintf("%#v %#v %#v", day, month, year))
		logTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		if logTime.After(start) && logTime.Before(end) || logTime == start || logTime == end {
			result = append(result, str)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("zero logs in diapazon nahui")
	}
	return result, nil
}
