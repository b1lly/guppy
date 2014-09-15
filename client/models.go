package main

import "log"

type GuppyLog struct{}

func (gl *GuppyLog) Error(v ...interface{}) {
	log.Println("[guppy]", v)
}

func (gl *GuppyLog) Info(v ...interface{}) {
	log.Println("[guppy]", v)
}

func (gl *GuppyLog) Warning(v ...interface{}) {
	log.Println("[guppy]", v)
}
