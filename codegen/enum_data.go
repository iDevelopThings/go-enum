package codegen

import (
	"bytes"
	"github.com/iancoleman/strcase"
	"strconv"
	"strings"
)

type EnumType string

const (
	EnumTypeStruct  EnumType = "struct"
	EnumTypeComment EnumType = "comment"
)

type Enum struct {
	DefinedName   string
	GeneratedName string

	ValueType string
	Values    []*EnumValue

	GeneratedFrom EnumType

	GenCamelCaseName      string
	GenLowerCamelCaseName string
	GenEnumListName       string
	GenEnumDefName        string
	GenInstanceVariable   string

	GeneratedFile bytes.Buffer
}

func NewEnum(definedName string) *Enum {

	generatedName := strings.Replace(definedName, "Enum", "", -1)
	if strings.HasSuffix(generatedName, "Def") {
		generatedName = strings.Replace(generatedName, "Def", "", -1)
	}

	e := &Enum{
		DefinedName:   definedName,
		GeneratedName: generatedName,
		Values:        make([]*EnumValue, 0),
	}

	e.GenCamelCaseName = strcase.ToCamel(generatedName)
	e.GenLowerCamelCaseName = strcase.ToLowerCamel(generatedName)

	e.GenEnumListName = pluralize.Plural(e.GenCamelCaseName)
	e.GenEnumDefName = pluralize.Singular(e.GenCamelCaseName)

	e.GenInstanceVariable = e.GenLowerCamelCaseName + "Instance"

	return e
}

func (e *Enum) AddValue(valueIdx int, name string) *EnumValue {
	value := NewEnumValue(name)
	value.Index = valueIdx

	e.Values = append(e.Values, value)
	return value
}

func (e *Enum) Finalize() {
	valueType := ""
	for _, value := range e.Values {
		value.Finalize()

		if valueType == "" {
			valueType = value.Type
		} else if valueType != value.Type {
			panic("Enum values have different types")
		}
	}

	e.ValueType = valueType
}

type EnumValue struct {

	// The index of this enum value in the list of values
	Index int

	// The name of the enum value as defined in the source code.
	DefinedName string

	// If there was a custom name defined on a struct enum via `name: "CustomName"`, then this will be set.
	CustomName string

	/**
	 * If the value of the enum struct uses strings, it will be set following this rule:
	 *     This will be a SCREAMING_SNAKE_CASE version of the defined name
	 *     If there was a custom name defined, it will be used instead.
	 * If the value of the enum struct uses ints
	 *     The value will start at 0 and increment by 1 for each enum value
	 */
	Value string

	// The value type defined on the enum struct
	Type string

	// If there was a description defined on a struct enum via `description: "This is a description"`, then this will be set.
	// If it was a comment, the description is the comment below the `// Name: ` defined comment
	Description string
}

// Finalize This is where we'll implement the changes to follow the logic for the fields/types
func (v *EnumValue) Finalize() {
	if v.Type == "string" {
		if v.CustomName != "" {
			v.Value = v.CustomName
		} else {
			v.Value = strcase.ToScreamingSnake(v.DefinedName)
		}
		// Surround the value with quotes so it can be used in the generated code
		v.Value = `"` + v.Value + `"`
	} else if v.Type == "int" {
		v.Value = strconv.Itoa(v.Index)
	}
}

func NewEnumValue(definedName string) *EnumValue {
	return &EnumValue{
		DefinedName: definedName,
		Type:        "string",
	}
}
