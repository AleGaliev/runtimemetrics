package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

const tmplStr = `{{ .Name }}
{{ .Fields }}
func ({{ .Pointer }} *{{.Name}}) Reset() {
{{range .Fields}}
	// {{.Name}}
	{{- if .IsSlice}}
	if {{ $.Pointer }}.{{.Name}} != nil {
		{{ $.Pointer }}.{{.Name}} = s.{{.Name}}[:0]
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
	if resetter, ok := s.{{.Name}}.(interface{ Reset() }); ok && {{ $.Pointer }}.{{.Name}} != nil {
		resetter.Reset()
	}
	{{- else}}
	{{ $.Pointer }}.{{.Name}} = {{.ZeroValue}}
	{{- end}}
{{end}}
}`

type FieldInfo struct {
	Name      string
	IsSlice   bool
	IsMap     bool
	IsPtr     bool
	IsStruct  bool
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
	//fmt.Println(structs)
	tmpl, err := template.New("resetMethod").Parse(tmplStr)
	for _, s := range structs {
		err = tmpl.Execute(os.Stdout, s)
		if err != nil {
			fmt.Println(err)

		}
	}
	//printAllStructs(structs)

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

//func convertInfoPackageInPackage(info []StructInfo) StructInfoPackage {
//	structsPackage := StructInfoPackage{
//		Package: make(map[string][]StructInfo),
//	}
//	for _, structs := range info {
//		fmt.Println(structs.Name)
//		structsPackage.Package[structs.Name] = append(structsPackage.Package[structs.Name], structs)
//	}
//	return structsPackage
//}

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

func printAllStructs(structs []StructInfo) {
	fmt.Printf("=== Найдено %d структур ===\n\n", len(structs))

	for i, s := range structs {
		fmt.Printf("%d. %s\n", i+1, s.Name)
		fmt.Printf("   Пакет:  %s\n", s.Package)
		fmt.Printf("   Комментарий:  %s\n", s.Comment)
		fmt.Printf("   Файл:   %s\n", filepath.Base(s.File))
		fmt.Printf("   Путь до файла:   %s\n", s.File)
		fmt.Printf("   Строка: %d\n", s.Line)

		if len(s.Fields) > 0 {
			fmt.Printf("   Поля (%d):\n", len(s.Fields))
			for _, f := range s.Fields {

				fmt.Printf("     %s %s\n",
					f.Name)
			}
		} else {
			fmt.Println("   Поля отсутствуют")
		}
		fmt.Println()
	}
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

func checkCommentGeneration(typeSpec *ast.TypeSpec, comments []*ast.CommentGroup) bool {
	if typeSpec.Doc != nil {
		for _, comment := range typeSpec.Doc.List {
			fmt.Println(comment.Text)
			if strings.HasPrefix(comment.Text, "generate:reset;") {
				return true
			}
		}
	}
	return false
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
