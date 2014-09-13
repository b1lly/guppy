package main

import (
	"fmt"
	"log"
)

type GuppyLog struct{}

func (gl *GuppyLog) Println(v interface{}) {
	log.Println(fmt.Sprintf("[guppy] %v", v))
}
