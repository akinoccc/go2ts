package core

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

// ParseDirRecursive 递归解析目录中的所有Go文件
func ParseDirRecursive(fset *token.FileSet, path string) (map[string]string, []*ast.File, map[string][]*ast.File, error) {
	allFiles := make([]*ast.File, 0)
	filesByPkg := make(map[string][]*ast.File)
	pkgPathByPkg := make(map[string]string)
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
			allFiles = append(allFiles, file)
			pkgName := file.Name.Name
			filesByPkg[pkgName] = append(filesByPkg[pkgName], file)
			if pkgPathByPkg[pkgName] == "" {
				pkgPathByPkg[pkgName] = filepath.Dir(path)
			}
		}
		return nil
	})
	return pkgPathByPkg, allFiles, filesByPkg, err
}

// ParseProject 解析整个项目并生成TypeScript接口定义
func ParseProject(path string) (string, error) {
	fset := token.NewFileSet()

	conf := types.Config{Importer: &LocalImporter{BaseDir: path}}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	_, allFiles, filesByPkg, err := ParseDirRecursive(fset, path)
	if err != nil {
		return "", err
	}
	if len(filesByPkg) == 0 {
		return "", fmt.Errorf("no files found")
	}

	// 类型检查所有文件
	//for _, files := range filesByPkg {
	//pkg := types.NewPackage(pkgPatyPkg[pkgName], pkgName)
	//check := types.NewChecker(&conf, fset, pkg, info)
	if _, err := conf.Check(path, fset, allFiles, info); err != nil {
		fmt.Println(11, err)
		return "", err
	}
	//}

	var tsInterfaces []string

	for ident, obj := range info.Defs {
		if obj, ok := obj.(*types.TypeName); ok {
			if structType, ok := obj.Type().Underlying().(*types.Struct); ok {
				tsInterfaces = append(tsInterfaces, GenerateTSInterface(ident.Name, structType))
			}
		}
	}

	return strings.Join(tsInterfaces, "\n"), nil
}
