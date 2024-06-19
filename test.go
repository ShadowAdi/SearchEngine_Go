package main

import (
	"fmt"
	"sync"
)

func test() {
	var wg sync.WaitGroup
	go func() {

		for i := 0; i < 10; i++ {
			fmt.Println("it is first", i)
		}
	}()

	counting()

	wg.Wait()

}
func counting() {
	for i := 0; i < 10; i++ {

		fmt.Println("it is second", i)
	}

}
