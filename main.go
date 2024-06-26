package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"regexp"
	"strings"
)

func goTypeToTSType(goType types.Type) string {
	// 如果 goType 是 *types.Named 类型，获取其基础类型
	var name string

	if named, ok := goType.(*types.Named); ok {
		name = named.Obj().Name()
		goType = named.Underlying()
	}

	switch t := goType.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return "boolean"
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64:
			return "number"
		case types.String:
			return "string"
		}
	case *types.Struct:
		switch name {
		case "Time":
			return "string"
		default:
			return name
		}
	case *types.Slice:
		return goTypeToTSType(t.Elem()) + "[]"
	}

	return "any"
}

func generateTSInterface(name string, t *types.Struct) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("interface %s {\n", name))

	for i := 0; i < t.NumFields(); i++ {
		field := t.Field(i)
		name := field.Name()

		re := regexp.MustCompile(`json:"([^"]+)"`)
		tag := t.Tag(i)
		match := re.FindStringSubmatch(tag)
		if len(match) >= 2 {
			name = match[1]
		}
		tsType := goTypeToTSType(field.Type())
		sb.WriteString(fmt.Sprintf("  %s: %s;\n", name, tsType))
	}

	sb.WriteString("}\n")
	return sb.String()
}

// 解析整个项目并生成TypeScript接口定义
func parseProject(path string) (string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.AllErrors)
	if err != nil {
		return "", err
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	var tsInterfaces []string

	for _, pkg := range pkgs {
		for p, file := range pkg.Files {
			astFile, err := parser.ParseFile(fset, p, nil, parser.AllErrors)
			if err != nil {
				return "", err
			}

			check := types.NewChecker(&conf, fset, types.NewPackage(path, file.Name.Name), info)
			if err := check.Files([]*ast.File{astFile}); err != nil {
				return "", err
			}

			for ident, obj := range info.Defs {
				if obj, ok := obj.(*types.TypeName); ok {
					if structType, ok := obj.Type().Underlying().(*types.Struct); ok {
						tsInterfaces = append(tsInterfaces, generateTSInterface(ident.Name, structType))
					}
				}
			}
		}
	}

	return strings.Join(tsInterfaces, "\n"), nil
}

func writeResultToTSFile(result string) error {
	tsFilePath := "/Users/akino/Desktop/Github/go2ts/example/ts-results/converted.d.ts"

	file, err := os.Create(tsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将 result 写入文件
	_, err = file.WriteString(result)
	if err != nil {
		return err
	}

	fmt.Printf("Result written to %s\n", tsFilePath)
	return nil
}

func main() {
	//if len(os.Args) < 2 {
	//    log.Fatalf("Usage: %s <path>", os.Args[0])
	//}

	path := "/Users/akino/Desktop/Github/go2ts/example"

	result, err := parseProject(path)
	if err != nil {
		log.Fatalf("Error parsing project: %v", err)
	}

	err = writeResultToTSFile(result)
	if err != nil {
		log.Fatalf("Error writing result to file: %v", err)
	}
}
