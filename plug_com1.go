package main

import "fmt"
import "dplug"
import "strconv"

func DoIt(p dplug.Parameters, r *dplug.Results) error {
	*r = dplug.Results{
		"woo": "aaaaa" + strconv.Itoa(p["num"].(int)) + "aaa",
	}
	return nil
}

func DoNothing(p dplug.Parameters, r *dplug.Results) error {
	*r = dplug.Results{
		"woo": "bbbb" + strconv.Itoa(p["num"].(int)) + "ccc",
	}
	return nil
}

func main() {
	gps := dplug.StartDPlugServer("test thingy")
	gps.RegisterMethod("doit", DoIt)
	gps.RegisterMethod("doit2it", DoNothing)
	err := gps.Serve()
	fmt.Println(err)
}
