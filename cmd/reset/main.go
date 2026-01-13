package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

const tmplStr = `package {{ .Package }}
func ({{ .Pointer }} *{{.Name}}) Reset() {
{{range .Fields}}
	{{- if .IsSlice}}
	if {{ $.Pointer }}.{{.Name}} != nil {
		{{ $.Pointer }}.{{.Name}} = {{ $.Pointer }}.{{.Name}}[:0]
	}
	{{- else if .IsString}}
	if {{ $.Pointer }}.{{.Name}} != "" {
		{{ $.Pointer }}.{{.Name}} = ""
	}
	{{- else if .IsInt}}
	if {{ $.Pointer }}.{{.Name}} != 0 {
		{{ $.Pointer }}.{{.Name}} = 0
	}
	{{- else if .IsMap}}
	if {{ $.Pointer }}.{{.Name}} != nil {
		clear({{ $.Pointer }}.{{.Name}})
	}
	{{- else if .IsPtr}}
	if {{ $.Pointer }}.{{.Name}} != nil {
		*{{ $.Pointer }}.{{.Name}} = {{.ZeroValue}}
	}
	{{- else if .IsStruct}}
	if resetter, ok := {{ $.Pointer }}.{{.Name}}.(interface{ Reset() }); ok && {{ $.Pointer }}.{{.Name}} != nil {
		resetter.Reset()
	}
	{{- else}}
	{{ $.Pointer }}.{{.Name}} = nil
	{{- end}}
{{- end }}
}
`

type FieldInfo struct {
	Name      string
	IsSlice   bool
	IsMap     bool
	IsPtr     bool
	IsStruct  bool
	IsString  bool
	IsInt     bool
	IsMutex   bool
	ZeroValue string
}

type StructInfo struct {
	Package string
	Pointer string
	File    string
	Name    string
	Comment string
	Line    int
	Fields  []FieldInfo
}

func main() {
	pwd := os.Getenv("PWD")
	projectDir := pwd

	structs, err := findAllStructsInProject(projectDir)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}
	tmpl, err := template.New("resetMethod").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
	structsList := convertInfoPackageInPackage(structs)

	for _, structs := range structsList {

		filePath := structs[0].File

		dir := filepath.Dir(filePath)

		newPath := filepath.Join(dir, "reset.gen.go")

		var buf bytes.Buffer
		for _, s := range structs {

			err = tmpl.Execute(&buf, s)

			if err != nil {
				log.Fatal(err)
			}

		}
		file, err := os.OpenFile(newPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.Write(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}

}

func findAllStructsInProject(rootDir string) ([]StructInfo, error) {
	var allStructs []StructInfo
	fset := token.NewFileSet()

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		// Только .go файлы
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Пропускаем тесты и моки
		if filepath.Ext(path) == ".go" && (strings.Contains(filepath.Base(path)[:len(filepath.Base(path))-3], "_test") || strings.Contains(filepath.Base(path)[:len(filepath.Base(path))-3], "_mocks")) {
			return nil
		}

		//Парсим файл
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)

		if err != nil {
			// Пропускаем файлы с ошибками
			return nil
		}

		//Находим структуры в файле
		structs := findStructsInFile(fset, f, path)
		allStructs = append(allStructs, structs...)
		return nil
	})

	return allStructs, err
}

func convertInfoPackageInPackage(info []StructInfo) map[string][]StructInfo {
	structsPackage := make(map[string][]StructInfo)

	for _, structs := range info {
		structsPackage[structs.Package] = append(structsPackage[structs.Package], structs)
	}
	return structsPackage
}

func findStructsInFile(fset *token.FileSet, f *ast.File, fileName string) []StructInfo {
	var structs []StructInfo
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Doc != nil {
			for _, comment := range genDecl.Doc.List {
				if strings.Contains(comment.Text, "generate:reset") {
					for _, spec := range genDecl.Specs {
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}
						structType, isStruct := typeSpec.Type.(*ast.StructType)
						if isStruct {
							runes := []rune(typeSpec.Name.Name)
							firstChar := unicode.ToLower(runes[0])
							info := StructInfo{
								Package: f.Name.Name,
								Pointer: string(firstChar),
								File:    fileName,
								Name:    typeSpec.Name.Name,
								Comment: comment.Text,
								Line:    fset.Position(comment.Pos()).Line,
								Fields:  []FieldInfo{},
							}

							if structType.Fields != nil {
								for _, field := range structType.Fields.List {
									fieldType := exprToString(field.Type)
									if len(field.Names) > 0 {
										for _, name := range field.Names {
											info.Fields = append(info.Fields, FieldInfo{
												Name:      name.Name,
												IsSlice:   strings.HasPrefix(fieldType, "[]"),
												IsMap:     strings.Contains(fieldType, "map["),
												IsPtr:     strings.HasPrefix(fieldType, "*"),
												IsStruct:  !strings.ContainsAny(fieldType, "[]*map"),
												IsString:  strings.HasPrefix(fieldType, "string"),
												IsInt:     strings.HasPrefix(fieldType, "int"),
												IsMutex:   strings.HasPrefix(fieldType, "sync.Mutex"),
												ZeroValue: getZeroValue(fieldType),
											})
										}
									}
								}
							}
							structs = append(structs, info)
						}
					}
				}
			}
		}
	}

	return structs
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s",
			exprToString(t.Key),
			exprToString(t.Value))
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func getZeroValue(t string) string {
	switch {
	case strings.HasPrefix(t, "*"), strings.HasPrefix(t, "[]"), strings.Contains(t, "map["):
		return "nil"
	case t == "int", t == "int8", t == "int16", t == "int32", t == "int64",
		t == "uint", t == "uint8", t == "uint16", t == "uint32", t == "uint64",
		t == "float32", t == "float64", t == "complex64", t == "complex128":
		return "0"
	case t == "string":
		return `""`
	case t == "bool":
		return "false"
	default:
		// Для структур и пользовательских типов
		return fmt.Sprintf("%s{}", t)
	}
}
