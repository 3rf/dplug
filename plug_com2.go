package main

import "fmt"
import "dplug"
import "log"
import "strconv"

func DoNothing(p dplug.Parameters, r *dplug.Results) error {
	log.Print(p)
	*r = dplug.Results{
		"woo": "bbbb" + strconv.Itoa(p["num"].(int)) + "ccc",
	}
	log.Print("returns: ", *r)
	return nil
}

func main() {
	gps := dplug.NewDPlugServer("test22")
	gps.RegisterMethod("doit2it", DoNothing)
	err := gps.Serve()
	fmt.Println(err)
}
