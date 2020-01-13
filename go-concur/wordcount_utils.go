package concur

import (
	"bufio"
	"os"
	"time"
)

var dir = "/Users/longfeixing/classics/"

type SlowIOReader struct{
	lines []string
	i int
}

func NewSlowIOReader(path string) *SlowIOReader {
	time.Sleep(time.Millisecond * 10)
	r := new(SlowIOReader)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	lines := make([]string, 0)
	for line, err := reader.ReadSlice('\n'); err == nil; line, err = reader.ReadSlice('\n'){
		lines = append(lines, string(line))
	}
	r.lines = lines
	return r
}

func (r *SlowIOReader) ReadLine() *string {
	if r.i >= len(r.lines) {
		return nil
	}
	s := r.lines[r.i]
	r.i = r.i + 1
	return &s
}

func mergeMap(maps []map[string]int) map[string]int {
	if len(maps) == 0 {
		return nil
	}
	if len(maps) == 1 {
		return maps[0]
	}
	m := maps[0]
	for _, mp := range maps[1:] {
		for word, count := range mp {
			m[word] = m[word] + count
		}
	}
	return m
}

