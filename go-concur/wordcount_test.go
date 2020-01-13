package concur

import (
	"os"
	"runtime"
	"testing"
)

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

func BenchmarkWordCountWithChannelOnly(b *testing.B) {
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
