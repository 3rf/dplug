package main

import (
    "fmt"
    "goplug"
)


func main() {
    goplug.TestInitSession() 
    r := goplug.Results{}
    err := goplug.CallPluginMethod("test","doit", goplug.Parameters{"num":154}, &r)
    fmt.Println(err)
    fmt.Println(r)
}

