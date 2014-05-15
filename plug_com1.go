package main

import "log"
import "net/http"
import "net"
import "net/rpc"

type Remote struct{}

func (r *Remote) AddSomeNumbers(nums *[]int, result *int) error {
    log.Print(*nums)
	sum := 0
	for _, num := range *nums {
		sum += num
	}
	*result = sum
    log.Print("returns: ", *result)
	return nil
}

func main() {
	rpc.Register(new(Remote))
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
