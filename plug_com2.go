package main

import "fmt"
import "dplug"
import "flag"
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
func DoNothing(p dplug.Parameters, r *dplug.Results) error {
	log.Print(p)
	*r = dplug.Results{
		"woo": "bbbb" + strconv.Itoa(p["num"].(int)) + "ccc",
	}
	log.Print("returns: ", *r)
	return nil
}

func main() {
	var port = flag.Int("dplugport", 0, "port to run plugin on")
	flag.Parse()

	gps := dplug.DPlugServer{
		Self: dplug.Plugin{
			Name: "test2",
			Port: *port},
		Methods: map[string]dplug.MethodHandler{},
	}
	gps.RegisterMethod("doit2it", dplug.MethodHandlerFunc(DoNothing))
	err := gps.Serve()
	fmt.Println(err)
}
