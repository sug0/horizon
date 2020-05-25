package script

import (
    "fmt"

    "github.com/d5/tengo/v2"
    "github.com/d5/tengo/v2/parser"
)

type Script struct {
    *tengo.Bytecode
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

func (s *Script) BootstrapVM(globals ...tengo.Object) *tengo.VM {
    return tengo.NewVM(s.Bytecode, globals, -1)
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
    c := tengo.NewCompiler(srcFile, symbols, nil, modules, nil)
    err = c.Compile(file)
    if err != nil {
        err = fmt.Errorf("script: failed to compile script: %w", err)
        return nil, err
    }
    return &Script{c.Bytecode()}, nil
}
