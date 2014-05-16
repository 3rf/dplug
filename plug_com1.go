package main

import "fmt"
import "dplug"
import "log"
import "strconv"

func DoIt(p dplug.Parameters, r *dplug.Results) error {
	log.Print(p)
	*r = dplug.Results{
		"woo": "aaaaa" + strconv.Itoa(p["num"].(int)) + "aaa",
	}
	log.Print("returns: ", *r)
	return nil
}

func main() {
	gps := dplug.DPlugServer{
		Self: dplug.Plugin{
			Name:        "test",
			MethodNames: []string{},
			Port:        1234},
		Methods: map[string]dplug.MethodHandler{},
	}
	gps.RegisterMethod("doit", dplug.MethodHandlerFunc(DoIt))
	err := gps.Serve()
	fmt.Println(err)
}
