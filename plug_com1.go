package main

import "goplug"
import "log"



func DoIt(p goplug.Parameters, r *goplug.Results) error {
    log.Print(p)
	*r = goplug.Results{
        "woo":"aaaaa" + p["num"].(string) + "aaa",

    }
    log.Print("returns: ", *r)
	return nil
}

func main() {
    gps := goplug.GoPlugServer{
        Self: goplug.Plugin{"test", []string{}, 1234},
        Methods: map[string]goplug.MethodHandler{},
    }
    gps.RegisterMethod("doit", goplug.MethodHandlerFunc(DoIt))
    gps.Serve()
}
