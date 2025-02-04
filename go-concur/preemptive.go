package concur

import "fmt"
import "runtime"
import "time"

func cpuIntensive(p *int) {
	for i := 1; i <= 100000000000; i++ {
		*p = i
		if i % 100000 == 0 {
			// uncomment the below to un-hang
			//println(i)
		}
	}
}

func preemptive() {
	runtime.GOMAXPROCS(1)

	x := 0
	go cpuIntensive(&x)

	time.Sleep(100 * time.Millisecond)

	// printed only after cpuIntensive is completely finished
	fmt.Printf("x = %d.\n", x)
}
