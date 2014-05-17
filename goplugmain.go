package main

import (
	"dplug"
	"fmt"
)

func main() {
	err := dplug.Initialize(
		dplug.Config{
			[]dplug.PluginRoute{
				{"./plug_com1", 1234},
				//{"./plug_com2", 11012},
			},
		})
	if err != nil {
		panic(err)
	}
	defer dplug.ShutDown()

	meth, _ := dplug.PluginMethods("test thingy")
	fmt.Println("METHODS:", meth)

	meth, _ = dplug.PluginsWithMethod("doit2it")
	fmt.Println("PLUGINS w/ 'doit2it':", meth)

	var r string
	err = dplug.CallPluginMethod("test thingy", "doit2it", 1664, &r)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
	}
    fmt.Println("Got", r)
}
