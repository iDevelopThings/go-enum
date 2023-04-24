package main

type StructErrorsEnumDef struct {
	PE_Regular_NoName string `name:"Special" description:"Expected type definition (type <name> { ... })"`
	// Expected type definition (type <name> { ... })
	PE_DocAbove string

	PE_DocRight string // Expected type definition (type <name> { ... })
}

// Start:ErrorsEnumDef

// Name: PE_Namespace_ExpectedIdentifier
// Expected identifier following 'namespace' keyword. e.g. 'namespace <name> { ... }'

// Name: PE_Namespace_ExpectedOpeningBracket
// Expected opening bracket following 'namespace <name> ', e.g. 'namespace <name> {'

// Name: PE_Namespace_ExpectedClosingBracket
// Expected closing bracket following 'namespace <name> { .... ', e.g. 'namespace <name> { .... }'

// End:ErrorsEnumDef

type NodeTypeEnum struct {
	// A test description for this element
	Unknown string `enum:"Unknown,UndefinedNodeType"`

	// A test description for this element
	TypeDefinition string `name:"TypeDefinition" description:"An ast node that represents a type definition, ie type Name { ... }"`
	TypeField      string `name:"Field" description:"An ast node that represents a field inside a type, ie 'Name: Type''"`
	TypeFieldType  string `name:"FieldType" description:"An ast node that represents the type defined on a field type, ie 'Name: Type''"`

	MemberAccessExpression string `name:"MemberAccessExpression" description:"An ast node that represents access to a member of a type, ie 'a.b', lhs -> rhs"`

	NamespaceReference string `name:"NamespaceReference" description:"An ast node that represents a reference to a namespace, ie lhs -> rhs(a.b) - lhs(a) is the namespace"`
	MemberReference    string `name:"MemberReference" description:"An ast node that represents a reference to a member of a type, ie lhs -> rhs(a.b) - rhs(b) is the member"`

	NativeID       string `name:"NativeID" description:"An ast node that represents a native id, ie 'nativeid: 123'"`
	NativeString   string `name:"NativeString" description:"An ast node that represents a native string, ie 'nativestring: 'abc''"`
	NativeInt      string `name:"NativeInt" description:"An ast node that represents a native int, ie 'nativeint: 123'"`
	NativeFloat    string `name:"NativeFloat" description:"An ast node that represents a native float, ie 'nativefloat: 123.123'"`
	NativeBool     string `name:"NativeBool" description:"An ast node that represents a native bool, ie 'nativebool: true'"`
	NativeDate     string `name:"NativeDate" description:"An ast node that represents a native date, ie 'nativedate: 2019-01-01'"`
	NativeDateTime string `name:"NativeDateTime" description:"An ast node that represents a native datetime, ie 'nativedatetime: 2019-01-01T12:00:00'"`
	NativeArray    string `name:"NativeArray" description:"An ast node that represents a native array, ie 'nativearray: [1,2,3]'"`
}
