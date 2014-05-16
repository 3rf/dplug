package main

import (
	"dplug"
	"fmt"
)

func main() {
	err := dplug.Init(dplug.Config{
		[]dplug.PluginRoute{
			{"./plug_com1", 1234},
		},
	})
	if err != nil {
		panic(err)
	}
	defer dplug.ShutDown()

	r := dplug.Results{}
	err = dplug.CallPluginMethod("test", "doit", dplug.Parameters{"num": 154}, &r)
	fmt.Println(err)
	fmt.Println(r["woo"])
}
