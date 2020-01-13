package concur

import "io/ioutil"

// concurrent with channel with select
// driver: BenchmarkWordCountWithChannelAndSelect
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
