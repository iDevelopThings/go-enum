package main

import (
	"fmt"
	"github.com/iDevelopThings/go-enum/codegen"
	"os"
	"testing"
)

func deleteGeneratedFiles() {
	os.Remove(codegen.Options.Output)
}

func Test_GeneratingFromComments(t *testing.T) {
	//defer deleteGeneratedFiles()

	codegen.Options = &codegen.GeneratorOptions{
		TypeNames: []string{"ErrorsEnumDef"},
		InputFile: "atest_enum_inputs.go",
		Output:    "atest_enum_inputs.gen.go",
		Template:  "codegen/templates/container.tmpl",
	}

	codegen.Generate()

	fmt.Println("done")
}

func Test_GeneratingFromStructs(t *testing.T) {
	//defer deleteGeneratedFiles()

	codegen.Options = &codegen.GeneratorOptions{
		TypeNames: []string{"StructErrorsEnumDef" /*, "NodeTypeEnum"*/},
		InputFile: "atest_enum_inputs.go",
		Output:    "atest_enum_inputs.gen.go",
		Template:  "codegen/templates/container.tmpl",
	}

	codegen.Generate()

	fmt.Println("done")
}

func Test_ContainerTemplate(t *testing.T) {
	//defer deleteGeneratedFiles()

	codegen.Options = &codegen.GeneratorOptions{
		TypeNames: []string{"StructErrorsEnumDef", "NodeTypeEnum"},
		InputFile: "atest_enum_inputs.go",
		Output:    "atest_enum_inputs.gen.go",
		Template:  "container",
	}

	codegen.Generate()

	fmt.Println("done")
}

/*func Test_GeneratedStructs(t *testing.T) {
	ev := StructErrors.PE_Regular_NoName
	fmt.Println("ev:", ev)

	ct := StructErrors

	testInput := func(enum StructError) {
		if enum != StructErrors.PE_Regular_NoName {
			panic("expected PE_Regular_NoName")
		}
		fmt.Println("enum:", enum)
	}
	testInput(StructErrors.PE_Regular_NoName)

	testAnyInput := func(enum any) {
		if _, ok := enum.(StructError); !ok {
			panic("expected StructError")
		}
		if enum != StructErrors.PE_Regular_NoName {
			panic("expected PE_Regular_NoName")
		}
		fmt.Println("enum:", enum)
	}
	testAnyInput(StructErrors.PE_Regular_NoName)

	fmt.Println("ct:", ct)

}*/
