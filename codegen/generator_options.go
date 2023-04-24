package codegen

import (
	_ "embed"
)

type GeneratorOptions struct {
	TypeNames []string
	InputFile string
	Output    string

	// A custom template path to use for generating the enum code
	Template string
}

var Options = &GeneratorOptions{
	Template: "container.tmpl",
}

func (o *GeneratorOptions) IsIncludedTypeName(name string) bool {
	for _, t := range o.TypeNames {
		if t == name {
			return true
		}
	}

	return false
}
