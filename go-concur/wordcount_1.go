package concur

import (
	"io/ioutil"
	"strings"
)

// driver: BenchmarkWordCountSeq

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

