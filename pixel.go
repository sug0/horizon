package horizon

import "github.com/d5/tengo/v2"

type builtinGetPixel struct {
    tengo.ObjectImpl
    width, height int64
    bitmap []byte
}


func (*builtinGetPixel) CanCall() bool {
    return true
}

func (gp *builtinGetPixel) Call(args ...tengo.Object) (tengo.Object, error) {
    r, g, b := &tengo.Int{}, &tengo.Int{}, &tengo.Int{}
    pixel := &tengo.Array{
        Value: []tengo.Object{r, g, b},
    }
    if len(args) < 2 {
        return pixel, nil
    }
    x, xok := args[0].(*tengo.Int)
    y, yok := args[1].(*tengo.Int)
    if xok && yok {
        x := x.Value
        y := y.Value
        if x < 0 || y < 0 || x >= gp.width || y >= gp.height {
            return pixel, nil
        }
        i := 3 * (y*gp.width + x)
        r.Value = int64(gp.bitmap[i+0])
        g.Value = int64(gp.bitmap[i+1])
        b.Value = int64(gp.bitmap[i+2])
    }
    return pixel, nil
}
