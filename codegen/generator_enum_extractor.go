package codegen

import (
	"github.com/dave/dst"
	"reflect"
	"regexp"
	"strings"
)

// We need a regex to check for `// Start:` and `// End:` comments.
// They need to ignore whitespace, be case-insensitive and be able to capture the name of the enum. (e.g. `// Start:ErrorsEnumDef`)
var startCommentRegex = regexp.MustCompile("(?i)//\\s*start:\\s*(.+)$")
var endCommentRegex = regexp.MustCompile("(?i)//\\s*end:\\s*(.+)$")

// We need some regex to check for `// Name:` and `// Value:` comments and extract the inner text.
var nameCommentRegex = regexp.MustCompile("(?i)//\\s*name:\\s*(.+)$")

//var valueCommentRegex = regexp.MustCompile("(?i)//\\s*value:\\s*(.+)$")

type EnumExtractor struct {
	Enums []*Enum
}

func NewEnumExtractor() *EnumExtractor {
	return &EnumExtractor{
		Enums: make([]*Enum, 0),
	}
}

func (f *EnumExtractor) IsEnumValueComment(comment string) (string, bool) {
	matches := nameCommentRegex.FindStringSubmatch(comment)
	if len(matches) == 0 {
		return "", false
	}

	return strings.TrimSpace(matches[1]), true
}

func (f *EnumExtractor) BuildEnumsFromScope(scope *dst.Scope) {

	for name, object := range scope.Objects {

		// Ensure it's a type
		if object.Kind != dst.Typ {
			continue
		}

		// Skip it if it's not included
		if !Options.IsIncludedTypeName(name) {
			continue
		}

		// We need to ensure that this object is a struct
		if structDec, ok := object.Decl.(*dst.TypeSpec); ok {
			if structType, ok := structDec.Type.(*dst.StructType); ok {
				f.buildEnumFromStruct(name, structDec, structType)
			}
		}
	}

}

func (f *EnumExtractor) buildEnumFromStruct(name string, structDec *dst.TypeSpec, structType *dst.StructType) {
	enum := NewEnum(name)
	enum.GeneratedFrom = EnumTypeStruct

	for i, field := range structType.Fields.List {

		elementName := field.Names[0].Name
		element := enum.AddValue(i, elementName)

		if typ, ok := field.Type.(*dst.Ident); ok {
			element.Type = typ.String()
		}

		if field.Tag != nil {
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

			if description, ok := tag.Lookup("description"); ok {
				element.Description = description
			}
			if name, ok := tag.Lookup("name"); ok {
				element.CustomName = name
			}
			if value, ok := tag.Lookup("value"); ok {
				element.Value = value
			}
		}

		// If no description has been set... we'll take a look at the comment
		if element.Description == "" {
			comments := field.Decorations()
			if len(comments.Start.All()) > 0 {
				element.Description = comments.Start.All()[0]
			}

			if element.Description != "" {
				element.Description = strings.TrimRight(element.Description, "\n")
			}
		}

	}

	enum.Finalize()
	f.Enums = append(f.Enums, enum)
}

func (f *EnumExtractor) BuildEnumsFromComments(comments []string) {
	// First we'll locate "// Start:ErrorsEnumDef" and "// End:ErrorsEnumDef" comments
	// We'll only take the comments in-between those two comments and build an enum from them.

	groups := make(map[string][]string)
	var currentGroup string
	for _, comment := range comments {

		if strings.TrimSpace(comment) == "" {
			continue
		}

		if matches := startCommentRegex.FindStringSubmatch(comment); len(matches) > 0 {
			currentGroup = strings.TrimSpace(matches[1])
			groups[currentGroup] = make([]string, 0)
		} else if matches := endCommentRegex.FindStringSubmatch(comment); len(matches) > 0 {
			currentGroup = ""
		} else if currentGroup != "" {
			groups[currentGroup] = append(groups[currentGroup], comment)
		}
	}

	// Now we'll iterate over the groups and build an enum from each group
	for groupName, groupComments := range groups {
		enum := NewEnum(groupName)
		enum.GeneratedFrom = EnumTypeComment

		enumIdx := 0
		for i, comment := range groupComments {
			if name, ok := f.IsEnumValueComment(comment); ok {
				current := enum.AddValue(enumIdx, name)
				enumIdx++

				// Now we'll look at the next comment to see if it's a description
				// We'll just check if the next comment is not a value comment and
				// make a bad assumption that this means it's a description....
				if i+1 < len(groupComments) {
					nextComment := groupComments[i+1]
					if _, ok := f.IsEnumValueComment(nextComment); !ok {
						desc := strings.TrimRight(nextComment, "\n")
						desc = strings.TrimLeft(desc, "\t")
						desc = strings.TrimLeft(desc, "//")
						desc = strings.TrimSpace(desc)
						current.Description = desc
					}
				}
			}
		}

		enum.Finalize()
		f.Enums = append(f.Enums, enum)
	}
}
