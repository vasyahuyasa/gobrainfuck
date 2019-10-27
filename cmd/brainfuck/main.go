package main

import (
	"io/ioutil"
	"log"
	"os"

	brainfuck "github.com/vasyahuyasa/gobrainfuck"
)

func main() {
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("can not read from stdin: ", err)
	}

	i := brainfuck.NewInterpreter()
	err = i.ParseString(string(text))
	if err != nil {
		log.Fatal("can not parse input programm: ", err)
	}

	err = i.Run()
	if err != nil {
		log.Fatal("error executing application: ", err)
	}

}
