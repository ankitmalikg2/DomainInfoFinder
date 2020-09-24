package main

// import (
// 	"fmt"
// 	"time"
// )

// const MAX = 50

// func main() {
// 	sem := make(chan int, MAX)
// 	count := 0
// 	for {
// 		sem <- 1 // will block if there is MAX ints in sem
// 		go func() {
// 			count++
// 			fmt.Println("hello again, world", count)

// 			time.Sleep(5 * time.Second)
// 			<-sem // removes an int from sem, allowing another to proceed
// 		}()
// 	}
// }
