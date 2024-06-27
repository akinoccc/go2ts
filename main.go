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
	"path/filepath"
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
		default:
			return "any"
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

// 递归解析目录中的所有Go文件
func parseDirRecursive(fset *token.FileSet, path string) (map[string][]*ast.File, error) {
	filesByPkg := make(map[string][]*ast.File)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				return err
			}
			pkgName := file.Name.Name
			filesByPkg[pkgName] = append(filesByPkg[pkgName], file)
		}
		return nil
	})
	return filesByPkg, err
}

// 解析整个项目并生成TypeScript接口定义
func parseProject(path string) (string, error) {
	fset := token.NewFileSet()

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	filesByPkg, err := parseDirRecursive(fset, path)
	if err != nil {
		return "", err
	}
	if len(filesByPkg) == 0 {
		return "", fmt.Errorf("no files found")
	}

	// 类型检查所有文件
	for pkgName, files := range filesByPkg {
		pkg := types.NewPackage(path, pkgName)
		check := types.NewChecker(&conf, fset, pkg, info)
		if err := check.Files(files); err != nil {
			return "", err
		}
	}

	var tsInterfaces []string

	for ident, obj := range info.Defs {
		if obj, ok := obj.(*types.TypeName); ok {
			if structType, ok := obj.Type().Underlying().(*types.Struct); ok {
				tsInterfaces = append(tsInterfaces, generateTSInterface(ident.Name, structType))
			}
		}
	}

	return strings.Join(tsInterfaces, "\n"), nil
}

func writeResultToTSFile(result string) error {
	tsFilePath := "./example/ts-results/converted.d.ts"

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

	path := "./"

	result, err := parseProject(path)
	if err != nil {
		log.Fatalf("Error parsing project: %v", err)
	}

	err = writeResultToTSFile(result)
	if err != nil {
		log.Fatalf("Error writing result to file: %v", err)
	}
}
