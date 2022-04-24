package main

const VERSION = "0.1.1-alpha.0"

func check(e any) {
	if e != nil {
		panic(e)
	}
}

func main() {
	Execute()
}
