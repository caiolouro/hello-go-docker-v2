package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Hello World with auto reload!")
		time.Sleep(time.Second * 3)
	}
}
