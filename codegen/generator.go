package codegen

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/fatih/color"
	pl "github.com/gertd/go-pluralize"
	"go/format"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFiles embed.FS

var pluralize = pl.NewClient()
var Generator = NewGenerator()

type CodeGenerator struct {
	buf bytes.Buffer // Accumulated output.

	PackageName string
	PackagePath string

	// The file that we're currently processing
	File *dst.File

	// All comments extracted from the above file
	FileComments []string

	// All the generated/extracted enum definitions
	Enums []*Enum
}

func NewGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (g *CodeGenerator) ProcessPackage() {
	cfg := &packages.Config{
		Mode:  packages.LoadSyntax,
		Tests: false,
		/*Logf: func(format string, args ...interface{}) {
			logger.Debugf(format, args...)
		},*/
	}

	pkgs, err := decorator.Load(cfg, Options.InputFile)
	if err != nil {
		panic(err)
	}

	pkg := pkgs[0]
	g.PackageName = pkg.Name
	g.PackagePath = pkg.GoFiles[0]
	g.File = pkg.Syntax[0]

	g.ExtractCommentsFromAst()

	if len(pkg.Syntax) != 1 {
		log.Fatalf("error: %d files found", len(pkg.Syntax))
	}
}

func (g *CodeGenerator) ExtractCommentsFromAst() {
	comments := make([]string, 0)

	commentMap := make(map[string]bool)

	for _, d := range g.File.Decls {
		comment := d.Decorations()
		for _, s := range comment.Start.All() {
			comments = append(comments, s)
			commentMap[s] = true
		}
	}

	decs := g.File.Decs.Name
	for _, s := range decs.All() {
		if _, ok := commentMap[s]; ok {
			continue
		}
		comments = append(comments, s)
	}

	g.FileComments = comments
}

func (g *CodeGenerator) BuildDefinitions() {
	enumExtractor := NewEnumExtractor()
	enumExtractor.BuildEnumsFromScope(g.File.Scope)
	enumExtractor.BuildEnumsFromComments(g.FileComments)

	if len(enumExtractor.Enums) == 0 {
		logger.Fatalf("no enums found in file %s", g.PackagePath)
		return
	}

	for _, enum := range enumExtractor.Enums {
		g.Enums = append(g.Enums, enum)
	}

	logger.Infof("Enums defined in source file:")

	g.PrintValues()

	// Filter the output so we only have the enums we want
	var enums []*Enum
	for _, enum := range g.Enums {
		for _, typeName := range Options.TypeNames {
			if enum.DefinedName == typeName {
				enums = append(enums, enum)
			}
		}
	}
	g.Enums = enums
}

func (g *CodeGenerator) GetTemplate(enum *Enum) *template.Template {
	base, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		logger.Fatal("instance template parse: ", err)
	}

	templateName := "templates/" + string(enum.GeneratedFrom) + ".tmpl"
	if Options.Template != "" {

		// Check first if our template name exists in the embedded files
		if _, err := templateFiles.ReadFile(Options.Template); err != nil && !os.IsNotExist(err) {
			// If it's not an os.IsNotExist error, we should panic
			if _, err = base.ParseFiles(Options.Template); err != nil {
				logger.Infof("Template names: %s", base.DefinedTemplates())
				logger.Fatal("instance template parse: ", err)
			}
		}

		templateName = Options.Template
	}

	return base.Lookup(filepath.Base(templateName))

	/*// Check if the template exists
	_, err = templateFiles.Open(templateName)
	// If it exists, we're not using a custom template or anything
	// Just early return
	if err == nil {

		t, err := template.ParseFS(templateFiles, templateName)
		if err != nil {
			logger.Fatal("instance template parse: ", err)
		}

		return t
	}

	// If it's not an os.IsNotExist error, we should panic
	if !os.IsNotExist(err) {
		logger.Fatal(err)
	}

	// Now we try to load the users custom template
	t, err := template.ParseFiles(Options.Template)
	if err != nil {
		logger.Fatal("instance template parse: ", err)
	}

	return t*/
}

func (g *CodeGenerator) WriteDefinitions() {
	for _, enum := range g.Enums {

		t := g.GetTemplate(enum)

		err := t.Execute(&enum.GeneratedFile, enum)
		if err != nil {
			logger.Fatal("Execute: ", err)
			return
		}
	}
}

func (g *CodeGenerator) WriteFiles() {
	g.Printf("// Generated by \"go_enum %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("// https://github.com/iDevelopThings/go-enum")
	g.Printf("\n")
	g.Printf("package %s", g.PackageName)
	g.Printf("\n")

	imports := []string{
		"github.com/iDevelopThings/go-enum/enum",
	}
	g.Printf("import (\n")
	for _, imp := range imports {
		g.Printf("\t . \"%s\"\n", imp)
	}
	g.Printf(")\n\n")

	for _, enum := range g.Enums {
		g.Printf("%s", enum.GeneratedFile.String())
	}

	src := g.format()

	err := os.WriteFile(Options.Output, src, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Generated enums written to %s\n", Options.Output)
}

func (g *CodeGenerator) PrintValues() {
	lightBold := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	gray := color.New(color.FgHiBlack, color.Bold).SprintFunc()
	blue := color.New(color.FgHiBlue).SprintFunc()

	maxNameLen := 0
	combNameLen := 0
	maxValueLen := 0
	for _, d := range g.Enums {
		for _, field := range d.Values {
			n := blue(field.DefinedName)
			if len(n) > maxNameLen {
				maxNameLen = len(n)
			}

			if field.CustomName != "" {
				cn := n + " [" + gray(field.CustomName) + "]   "
				if len(cn) > combNameLen {
					combNameLen = len(cn)
				}
			}

			if len(field.Value) > maxValueLen {
				maxValueLen = len(field.Value)
			}
		}
	}

	for _, enum := range g.Enums {
		fmt.Printf("Enum(%s): %s - %s\n", gray(enum.DefinedName), lightBold(enum.GeneratedName), gray("Defined as: ", lightBold(enum.GeneratedFrom)))
		for _, field := range enum.Values {
			formatString := fmt.Sprintf("   %%-%ds %%-%ds %s", maxNameLen, maxValueLen, field.Value)
			fmt.Printf(formatString, blue(field.DefinedName), gray(field.CustomName))
			fmt.Printf("\n")
		}
	}
}

func (g *CodeGenerator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

func (g *CodeGenerator) Printf(format string, args ...interface{}) {
	_, err := fmt.Fprintf(&g.buf, format, args...)
	if err != nil {
		logger.Fatal(err)
	}
}
