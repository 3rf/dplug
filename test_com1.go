package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := []int{1, 2, 3}
	var reply int
	err = client.Call("Remote.AddSomeNumbers", &args, &reply)
	if err != nil {
		log.Fatal("remote error:", err)
	}
	fmt.Println(reply)
}
