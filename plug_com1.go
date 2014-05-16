package main

import "fmt"
import "goplug"
import "log"
import "strconv"

func DoIt(p goplug.Parameters, r *goplug.Results) error {
	log.Print(p)
	*r = goplug.Results{
		"woo": "aaaaa" + strconv.Itoa(p["num"].(int)) + "aaa",
	}
	log.Print("returns: ", *r)
	return nil
}

func main() {
	gps := goplug.GoPlugServer{
		Self:    goplug.Plugin{"test", []string{}, 1234},
		Methods: map[string]goplug.MethodHandler{},
	}
	gps.RegisterMethod("doit", goplug.MethodHandlerFunc(DoIt))
	err := gps.Serve()
	fmt.Println(err)

}
