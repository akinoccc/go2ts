package main

import (
	"go2ts/core"
	"log"
)

func main() {
	//if len(os.Args) < 2 {
	//    log.Fatalf("Usage: %s <path>", os.Args[0])
	//}

	path := "./example"

	result, err := core.ParseProject(path)
	if err != nil {
		log.Fatalf("Error parsing project: %v", err)
	}

	err = core.WriteResultToTSFile(result)
	if err != nil {
		log.Fatalf("Error writing result to file: %v", err)
	}
}

//func main() {
//	cfg := &packages.Config{
//		Mode: packages.LoadSyntax,
//		Dir:  "/Users/akino/Desktop/Github/go2ts",
//	}
//
//	pkgs, err := packages.Load(cfg, "example/example2")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fset := token.NewFileSet()
//	conf := types.Config{Importer: packages.(fset, pkgs)}
//	info := &types.Info{
//		Uses: make(map[*ast.Ident]types.Object),
//	}
//
//	for _, pkg := range pkgs {
//		for _, file := range pkg.Syntax {
//			_, err := conf.Check(pkg.PkgPath, fset, []*ast.File{file}, info)
//			if err != nil {
//				fmt.Println(err)
//			}
//		}
//	}
//}
