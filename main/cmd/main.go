package main

import (
    "os"
    "unsafe"
    "strconv"
    "io/ioutil"

    "github.com/d5/tengo/v2"
    "github.com/sug0/lilliput"
    "github.com/sug0/horizon/script"
)

func main() {
    s, err := script.Compile([]byte(`pixel[0] += 5`))
    if err != nil {
        panic(err)
    }
    bufdata, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        panic(err)
    }
    decoder, err := lilliput.NewDecoder(bufdata)
    if err != nil {
        panic(err)
    }
    defer decoder.Close()
    header, err := decoder.Header()
    if err != nil {
        panic(err)
    }
    pixtyp := header.PixelType()
    if pixtyp.Depth() != 8 {
        panic("unsupported depth: " + strconv.Itoa(pixtyp.Depth()))
    }
    if pixtyp.Channels() != 3 {
        panic("unsupported channels: " + strconv.Itoa(pixtyp.Channels()))
    }
    framebuf := lilliput.NewFramebuffer(header.Width(), header.Height())
    defer framebuf.Clear()
    bitmap := ((*struct { Buf []byte })(unsafe.Pointer(framebuf))).Buf
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
        setPixel(&pixel, 0, bitmap[i+0])
        setPixel(&pixel, 1, bitmap[i+1])
        setPixel(&pixel, 2, bitmap[i+2])

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
