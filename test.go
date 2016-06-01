package scatter

import (
	"fmt"
	"time"
	"math/rand"
)

func TestRoutineFunc (goroutine_index int, name interface{}) {
	var tmp *string
	tmp = name.(*string)
	fmt.Printf("my test info : %s : %d\n", *tmp, goroutine_index)

	time.Sleep(200 * time.Millisecond + time.Duration(rand.Intn(200)) * time.Millisecond)
}

func Test() {

	// create a new test instance
	scatter_instance := NewScatter()

	var TestTimes int
	TestTimes = 1000
	var name string
	scatter_instance.Init()
	scatter_instance.Do(TestRoutineFunc, TestTimes, &name)
}
