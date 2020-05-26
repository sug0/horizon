package script

import (
    "fmt"

    "github.com/d5/tengo/v2"
    "github.com/d5/tengo/v2/parser"
)

type Script struct {
    *tengo.Bytecode

    symbols *tengo.SymbolTable
    modules *tengo.ModuleMap
}

func (s *Script) BootstrapVM(globals ...tengo.Object) *tengo.VM {
    actualGlobals := make([]tengo.Object, 32)
    for i := 0; i < len(globals); i++ {
        actualGlobals[i] = globals[i]
    }
    return tengo.NewVM(s.Bytecode, actualGlobals, -1)
}

func Compile(source []byte) (*Script, error) {
    fileSet := parser.NewFileSet()
    srcFile := fileSet.AddFile("(main)", -1, len(source))
    p := parser.NewParser(srcFile, source, nil)
    file, err := p.ParseFile()
    if err != nil {
        err = fmt.Errorf("script: failed to parse script: %w", err)
        return nil, err
    }
    symbols := buildSymTab()
    modules := tengo.NewModuleMap()
    c := tengo.NewCompiler(srcFile, symbols, nil, modules, nil)
    err = c.Compile(file)
    if err != nil {
        err = fmt.Errorf("script: failed to compile script: %w", err)
        return nil, err
    }
    s := &Script{
        Bytecode: c.Bytecode(),
        modules: modules,
        symbols: symbols,
    }
    return s, nil
}

func buildSymTab() (symbols *tengo.SymbolTable) {
    symbols = tengo.NewSymbolTable()

    // globals[0] --> output
    symbols.Define("output")

    // globals[1] --> x
    symbols.Define("x")

    // globals[2] --> y
    symbols.Define("y")

    // globals[3] --> i64
    symbols.Define("i64")

    // globals[4] --> f64
    symbols.Define("f64")

    // globals[5] --> get_pixel
    symbols.Define("get_pixel")

    // globals[6] --> width
    symbols.Define("width")

    // globals[7] --> height
    symbols.Define("height")

    // globals[8] --> ctx
    symbols.Define("ctx")

    return
}
