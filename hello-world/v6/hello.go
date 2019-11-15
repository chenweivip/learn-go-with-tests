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

	if language == chinese {
		return chineseHelloPrefix + name
	}

	if language == french {
		return frenchHelloPrefix + name
	}

	return englishHelloPrefix + name
}

func main() {
	fmt.Println(Hello("world", ""))
}
