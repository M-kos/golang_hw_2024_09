package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("not enough args")
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
