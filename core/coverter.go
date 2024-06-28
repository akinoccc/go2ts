package core

import (
	"fmt"
	"go/types"
	"os"
	"regexp"
	"strings"
)

func GoTypeToTSType(goType types.Type) string {
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
		return GoTypeToTSType(t.Elem()) + "[]"
	}

	return "any"
}

func GenerateTSInterface(name string, t *types.Struct) string {
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
		tsType := GoTypeToTSType(field.Type())
		sb.WriteString(fmt.Sprintf("  %s: %s;\n", name, tsType))
	}

	sb.WriteString("}\n")
	return sb.String()
}

func WriteResultToTSFile(result string) error {
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
