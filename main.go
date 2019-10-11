package main

import "github.com/janabe/cscoupler/server"

func main() {
	server := server.NewServer()
	server.Run()
}
