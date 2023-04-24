package enum

type EnumContainer[TEnumType any, TEnumValueType comparable] struct {
	Values []TEnumType
	Map    map[TEnumValueType]TEnumType
}

type Enum[T any] struct {
	Name        string
	Description string
	Value       T
}

func (e Enum[T]) String() string {
	return e.Name
}
