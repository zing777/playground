package concur

import (
	"io/ioutil"
	"strings"
	"sync"
)

// driver:
// BenchmarkWordCountSharedDataOneCore
// BenchmarkWordCountSharedDataFourCore
// BenchmarkWordCountSharedDataAllCore
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

