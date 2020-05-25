package main

import (
    "os"
    "io/ioutil"

    "gocv.io/x/gocv"
    "github.com/sug0/horizon"
    "github.com/sug0/horizon/script"
)

const glitchProgram = `
max := func(x, y) {
    if x > y {
        return x
    }
    return y
}

// check if we saved last values
if !ctx.r {
    ctx.r = pixel[0]
    ctx.g = pixel[1]
    ctx.b = pixel[2]
}

pixel[0] = max(pixel[0], pixel[1])
pixel[1] = max(pixel[1], pixel[2])
pixel[2] = max(pixel[2], pixel[0])

pixel[0] *= ctx.r
pixel[1] *= ctx.g
pixel[2] = (ctx.b - (ctx.r|ctx.g)) * pixel[2]
`

func main() {
    s, err := script.Compile([]byte(glitchProgram))
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
    newmat, err := horizon.Glitch(s, mat)
    if err != nil {
        panic(err)
    }
    defer newmat.Close()
    bufdata, err = gocv.IMEncode(gocv.PNGFileExt, newmat)
    if err != nil {
        panic(err)
    }
    os.Stdout.Write(bufdata)
}
