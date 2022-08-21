package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeWheel(t *testing.T) {
	var tw *TimeWheel[int]
	cb := func(id int) {
		if id == 100 {
			fmt.Println("===", id, time.Now().UnixNano()/1e6)
		}
		tw.Timeout(time.Second, id)
	}
	tw = NewTimeWheel(time.Microsecond*10000, cb)
	for i := 0; i < 101; i++ {
		//ch := make(chan bool, 1)
		//id := i
		//go func(id int) {
		//	for {
		//		<-ch
		//		//fmt.Println(time.Now().UnixNano() / 1e6)
		//		tw.Timeout(time.Second*1, id)
		//	}
		//}(id)

		//time.Sleep(1 * time.Second)
		tw.Timeout(time.Second*1, i)
	}
	time.Sleep(time.Hour * 24)
}

func TestTicker(t *testing.T) {
	fmt.Println("start:", time.Now().UnixNano()/1e6)
	for {
		time.Sleep(time.Second * 1)
		fmt.Println("timeout:", time.Now().UnixNano()/1e6)
	}
}
