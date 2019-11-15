package main

import "fmt"

const chinese = "Chinese"
const french = "French"
const englishHelloPrefix = "Hello, "
const chineseHelloPrefix = "你好, "
const frenchHelloPrefix = "Bonjour, "

// Hello returns a personalised greeting in a given language
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	prefix := englishHelloPrefix

	switch language {
	case french:
		prefix = frenchHelloPrefix
	case chinese:
		prefix = chineseHelloPrefix
	}

	return prefix + name
}

func main() {
	fmt.Println(Hello("world", ""))
}
