package main

import (
    "os"
    "fmt"
    "io/ioutil"

    "gocv.io/x/gocv"
    "github.com/d5/tengo/v2"
    "github.com/sug0/horizon/script"
)

func main() {
    s, err := script.Compile([]byte(`pixel[0] += 5; pixel[1] += 10`))
    if err != nil {
        panic(err)
    }
    bufdata, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        panic(err)
    }
    mat, err := gocv.IMDecode(bufdata, gocv.IMReadColor)
    if err != nil {
        panic(err)
    }
    defer mat.Close()
    bitmap := mat.DataPtrUint8()
    pixel := tengo.Array{
        Value: []tengo.Object{
            &tengo.Int{},
            &tengo.Int{},
            &tengo.Int{},
        },
    }
    vm := s.BootstrapVM(&pixel)
    for i := 0; i < len(bitmap); i += 3 {
        // set pixel value
        setPixel(&pixel, 0, int64(bitmap[i+0]))
        setPixel(&pixel, 1, int64(bitmap[i+1]))
        setPixel(&pixel, 2, int64(bitmap[i+2]))

        // print before
        fmt.Printf("%v --> ", pixel.Value)

        // run vm
        err = vm.Run()
        if err != nil {
            panic(err)
        }

        // print after
        fmt.Printf("%v\n", pixel.Value)
    }
}

func setPixel(pixel *tengo.Array, i int, x int64) {
    v := pixel.Value[i].(*tengo.Int)
    v.Value = x
}
