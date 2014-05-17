package main

import "fmt"
import "dplug"
import "strconv"

func DoIt(p int, r *string) error {
	*r = "aaaaa" + strconv.Itoa(p) + "aaa"
	return nil
}

func DoNothing(p int, r *string) error {
	*r = "bbbb" + strconv.Itoa(p) + "ccc"
	return nil
}

func main() {
	gps := dplug.NewDPlugServer("test thingy")
	gps.RegisterMethod("doit", DoIt)
	gps.RegisterMethod("doit2it", DoNothing)
	err := gps.Serve()
	fmt.Println(err)
}
