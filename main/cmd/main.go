package main

import (
    "os"
    "flag"
    "time"
    "context"
    "io/ioutil"

    "gocv.io/x/gocv"
    "github.com/sug0/horizon"
    "github.com/sug0/horizon/script"
)

func main() {
    var imagePath string
    var outPath string
    var glitchProgramPath string

    flag.StringVar(&imagePath, "in", "-", "The input image path. If empty, reads from stdin.")
    flag.StringVar(&outPath, "out", "-", "The output image path. If empty, writes to stdout.")
    flag.StringVar(&glitchProgramPath, "script", "", "The glitching script path.")
    flag.Parse()

    if glitchProgramPath == "" {
        panic("github.com/sug0/horizon/main/cmd: you need to specify a glitch script")
    }
    glitchProgram, err := ioutil.ReadFile(glitchProgramPath)
    if err != nil {
        panic(err)
    }
    s, err := script.Compile(glitchProgram)
    if err != nil {
        panic(err)
    }

    var bufdata []byte
    if imagePath == "-" {
        bufdata, err = ioutil.ReadAll(os.Stdin)
        if err != nil {
            panic(err)
        }
    } else {
        bufdata, err = ioutil.ReadFile(imagePath)
        if err != nil {
            panic(err)
        }
    }
    mat, err := gocv.IMDecode(bufdata, gocv.IMReadColor)
    if err != nil {
        panic(err)
    }
    defer mat.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
    defer cancel()
    newmat, err := horizon.Glitch(ctx, s, mat)
    if err != nil {
        panic(err)
    }
    defer newmat.Close()
    bufdata, err = gocv.IMEncode(gocv.PNGFileExt, newmat)
    if err != nil {
        panic(err)
    }
    if outPath == "-" {
        _, err = os.Stdout.Write(bufdata)
        if err != nil {
            panic(err)
        }
        return
    }
    err = ioutil.WriteFile(outPath, bufdata, 0644)
    if err != nil {
        panic(err)
    }
}
