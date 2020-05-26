package horizon

import "github.com/d5/tengo/v2"

type builtinConvF64 struct {
    tengo.ObjectImpl
}

type builtinConvI64 struct {
    tengo.ObjectImpl
}

func (*builtinConvF64) CanCall() bool {
    return true
}

func (*builtinConvF64) Call(args ...tengo.Object) (tengo.Object, error) {
    f := &tengo.Float{}
    if len(args) < 1 {
        return f, nil
    }
    if x, ok := args[0].(*tengo.Int); ok {
        f.Value = float64(x.Value)
    }
    return f, nil
}

func (*builtinConvI64) CanCall() bool {
    return true
}

func (*builtinConvI64) Call(args ...tengo.Object) (tengo.Object, error) {
    i := &tengo.Int{}
    if len(args) < 1 {
        return i, nil
    }
    if x, ok := args[0].(*tengo.Float); ok {
        i.Value = int64(x.Value)
    }
    return i, nil
}
