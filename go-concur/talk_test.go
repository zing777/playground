package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
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

// seq

func readFileSeq(path string, countMap map[string]int) map[string]int {
	reader := NewSlowIOReader(path)

	for line := reader.ReadLine(); line != nil; line = reader.ReadLine(){
		words := strings.Split(*line, " ")
		for _, word := range words {
			countMap[word] = countMap[word] + 1
		}
	}
	return countMap
}

func wordCountSequential() {
	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	maps := make([]map[string]int, 0, len(paths))
	var countMap = make(map[string]int)
	for _, path := range paths {
		maps = append(maps, readFileSeq(dir + path.Name(), countMap))
	}
	mergeMap(maps)
}

// concurrent with map

func readFileWithMap(path string, countMap map[string]int, lock *sync.Mutex, wg *sync.WaitGroup) {
	reader := NewSlowIOReader(path)
	for line := reader.ReadLine(); line != nil; line = reader.ReadLine(){
		words := strings.Split(*line, " ")
		for _, word := range words {
			lock.Lock()
			countMap[word] = countMap[word] + 1
			lock.Unlock()
		}
	}
	wg.Done()
}

func wordCountSharedData() {
	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var mapWG sync.WaitGroup
	var lock sync.Mutex
	var countMap = make(map[string]int)
	for _, path := range paths {
		mapWG.Add(1)
		go readFileWithMap(dir + path.Name(), countMap, &lock, &mapWG)
	}
	mapWG.Wait()
}

// concurrent with channel

func readFileWithChannel(path string, c chan map[string]int) {
	countMap := make(map[string]int)
	reader := NewSlowIOReader(path)
	
	for line := reader.ReadLine(); line != nil; line = reader.ReadLine(){
		words := strings.Split(*line, " ")
		for _, word := range words {
			countMap[word] = countMap[word] + 1
		}
	}
	c <- countMap
}


func wordCountWithChannel() {
	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	chans := make([]chan map[string]int, len(paths))
	for i, path := range paths {
		chans[i] = make(chan map[string]int)
		go readFileWithChannel(dir + path.Name(), chans[i])
	}
	maps := make([]map[string]int, 0, len(paths))
	for i := 0; i < len(paths); i ++ {
		maps = append(maps, <- chans[i])
	}
	mergeMap(maps)
}

// concurrent with channel with select

func wordCountWithChannelAnsSelect() {
	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	chann := make(chan map[string]int, len(paths))
	for _, path := range paths {
		go readFileWithChannel(dir + path.Name(), chann)
	}
	mapp := make(map[string]int)
	for i := 0; i < len(paths); i ++ {
		select {
		// fan out
		case mapp2 := <- chann:
			mapp = mergeMap([]map[string]int{mapp, mapp2})
		}
	}
}



// benchmark drivers

func BenchmarkWordCountSeq(b *testing.B)  {
	for i := 0; i < b.N; i ++ {
		wordCountSequential()
	}
}

func BenchmarkWordCountSharedDataOneCore(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i ++ {
		wordCountSharedData()
	}
}

func BenchmarkWordCountSharedDataFourCore(b *testing.B) {
	runtime.GOMAXPROCS(4)
	for i := 0; i < b.N; i ++ {
		wordCountSharedData()
	}
}

func BenchmarkWordCountSharedDataAllCore(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < b.N; i ++ {
		wordCountSharedData()
	}
}

func BenchmarkWordCountWithChannel(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		wordCountWithChannel()
	}
}

func BenchmarkWordCountWithChannelAndSelect(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		wordCountWithChannelAnsSelect()
	}
}

func TestMain(m *testing.M) {
	println("cpu count: ", runtime.NumCPU())
	os.Exit(m.Run())
}
