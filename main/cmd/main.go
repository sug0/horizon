package main

import (
    "fmt"

    "github.com/d5/tengo/v2"
    "github.com/d5/tengo/v2/parser"
)

type script struct {
    bytecode *tengo.Bytecode
}

var (
    symbols *tengo.SymbolTable
    modules *tengo.ModuleMap
)

func init() {
    symbols = tengo.NewSymbolTable()
    modules = tengo.NewModuleMap()

    // globals[0] --> pixel
    symbols.Define("pixel")
}

func main() {
    s, err := compileScript([]byte(`pixel[0] += 5`))
    if err != nil {
        panic(err)
    }
    pixel := tengo.Array{
        Value: []tengo.Object{
            &tengo.Int{},
            &tengo.Int{},
            &tengo.Int{},
        },
    }
    vm := s.bootstrapVM(&pixel)
    bitmap := []int64{
        5, 4, 1,
        2, 1, 0,
        0, 4, 2,
    }
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

func (s script) bootstrapVM(pixel *tengo.Array) *tengo.VM {
    globals := []tengo.Object{pixel}
    return tengo.NewVM(s.bytecode, globals, -1)
}

func compileScript(source []byte) (s script, err error) {
    var file *parser.File
    fileSet := parser.NewFileSet()
    srcFile := fileSet.AddFile("(main)", -1, len(source))
    p := parser.NewParser(srcFile, source, nil)
    file, err = p.ParseFile()
    if err != nil {
        return
    }
    c := tengo.NewCompiler(srcFile, symbols, nil, modules, nil)
    err = c.Compile(file)
    if err != nil {
        return
    }
    s.bytecode = c.Bytecode()
    return
}
