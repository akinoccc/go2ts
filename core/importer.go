package core

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
)

type LocalImporter struct {
	BaseDir string
}

func (li LocalImporter) Import(path string) (*types.Package, error) {
	config := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
		Dir:  li.BaseDir,
		Env:  os.Environ(),
	}

	pkgs, err := packages.Load(config, path)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found for path: %s", path)
	}
	if len(pkgs[0].Errors) > 0 {
		return nil, pkgs[0].Errors[0]
	}

	pkg := pkgs[0]
	fmt.Println(11, pkg.PkgPath)

	return types.NewPackage(pkg.PkgPath, pkg.Name), nil
}
