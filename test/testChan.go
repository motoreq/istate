package main

import (
	"fmt"
)

func main() {
	a := make(map[int]struct{})
	c := make(chan struct{})

	go func(a1 map[int]struct{}, c chan struct{}) {
		for i := 0; i < 100; i++ {
			a1[i] = struct{}{}
		}
		c <- struct{}{}
	}(a, c)

	go func(a2 map[int]struct{}, c chan struct{}) {
		for i := 100; i < 200; i++ {
			a2[i] = struct{}{}
		}
		c <- struct{}{}
	}(a, c)

	_, _ = <-c, <-c
	fmt.Println(len(a))
}
