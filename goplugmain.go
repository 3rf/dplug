package main

import (
	"dplug"
	"fmt"
)

func main() {
	err := dplug.Init(
		dplug.Config{
			[]dplug.PluginRoute{
				{"./plug_com1", 10101},
				{"./plug_com2", 1101},
			},
		})
	if err != nil {
		panic(err)
	}
	defer dplug.ShutDown()

	meth, _ := dplug.PluginMethods("test")
	fmt.Println("METHODS:", meth)

	meth, _ = dplug.PluginsWithMethod("doit2it")
	fmt.Println("PLUGINS w/ 'doit2it':", meth)

	r := dplug.Results{}
	err = dplug.CallPluginMethod("test", "doit2it", dplug.Parameters{"num": 154}, &r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r["woo"])
}
