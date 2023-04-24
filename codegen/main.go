package codegen

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/withmandala/go-log"
	"os"
	"strings"
)

var logger = log.New(os.Stdout).WithColor().WithDebug()

var (
	typeNames = flag.String("types", "", "comma-separated list of type names")
	output    = flag.String("output", "", "output file name; default <src dir>/enum_output.go")
	input     = flag.String("input", "", "input file name; default <src dir>/enum_location.go")

	templateName = flag.String("template", "", "a custom template path to use for generating the enum code")
)

func RunCli() {
	flag.Usage = usage
	flag.Parse()
	if !validateFlags() {
		flag.Usage()
		os.Exit(2)
	}

	Generate()
}

func validateFlags() bool {
	if *input == "" {
		logger.Error("Must provide an input file. This is the file where your enums are defined")
		return false
	}

	Options.InputFile = *input

	if *output == "" {
		logger.Error("Must provide an output file. This is the file where your enums will be generated")
		return false
	}

	Options.Output = *output

	if *typeNames == "" {
		logger.Warnf("No type names provided. Enums will be generated for all inside the provided input file: %s", *input)
	}

	Options.TypeNames = strings.Split(*typeNames, ",")

	if *templateName != "" {
		Options.Template = *templateName
	}

	return true
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of idevelopthings/go-enum\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-enum -input=myfile.go -output=myfile.gen.go -types=SomeEnumStruct\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-enum -input=myfile.go -output=myfile.gen.go -types=SomeEnumStruct -template=<some local go .tmpl file path>\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func Generate() {
	color.NoColor = false

	Generator.ProcessPackage()
	Generator.BuildDefinitions()
	Generator.WriteDefinitions()

	Generator.WriteFiles()
}
