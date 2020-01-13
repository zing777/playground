package concur

import (
	"io/ioutil"
	"strings"
)

// concurrent with channel
// driver: BenchmarkWordCountWithChannelOnly
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
