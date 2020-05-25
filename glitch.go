package horizon

import (
    "fmt"

    "gocv.io/x/gocv"
    "github.com/d5/tengo/v2"
    "github.com/sug0/horizon/script"
)

var (
    ErrNotRGBImage = fmt.Errorf("horizon: not an rgb image")
)

func Glitch(s *script.Script, mat gocv.Mat) (glitched gocv.Mat, err error) {
    if mat.Type() != gocv.MatTypeCV8UC3 {
        err = ErrNotRGBImage
        return
    }
    width := int64(mat.Cols())
    glitched = mat.Clone()
    bitmapMat := mat.DataPtrUint8()
    bitmapGlitched := glitched.DataPtrUint8()
    var x, y int64
    var coords [2]tengo.Int
    pixel := tengo.Array{
        Value: []tengo.Object{
            &tengo.Int{},
            &tengo.Int{},
            &tengo.Int{},
        },
    }
    persistent := tengo.Map{
        Value: make(map[string]tengo.Object),
    }
    vm := s.BootstrapVM(
        &pixel,
        &coords[0],
        &coords[1],
        &convI64{},
        &convF64{},
        &persistent,
    )
    for i := 0; i < len(bitmapMat); i += 3 {
        // set pixel value
        setPixel(&pixel, 0, int64(bitmapMat[i+0]))
        setPixel(&pixel, 1, int64(bitmapMat[i+1]))
        setPixel(&pixel, 2, int64(bitmapMat[i+2]))

        // run vm
        err = vm.Run()
        if err != nil {
            glitched.Close()
            err = fmt.Errorf("horizon: vm failed whilst running: %w", err)
            return
        }

        // get glitched pixel and set it
        bitmapGlitched[i+0] = byte(getPixel(&pixel, 0) & 0xff)
        bitmapGlitched[i+1] = byte(getPixel(&pixel, 1) & 0xff)
        bitmapGlitched[i+2] = byte(getPixel(&pixel, 2) & 0xff)

        // update coords
        if x < width {
            x++
        } else {
            x = 0
            y++
        }
        coords[0].Value = x
        coords[1].Value = y
    }
    return
}

func setPixel(pixel *tengo.Array, i int, x int64) {
    v := pixel.Value[i].(*tengo.Int)
    v.Value = x
}

func getPixel(pixel *tengo.Array, i int) int64 {
    v := pixel.Value[i].(*tengo.Int)
    return v.Value
}
