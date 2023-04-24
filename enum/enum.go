package enum

type EnumContainer[TEnumType any, TEnumValueType comparable] struct {
	Values []TEnumType
	Map    map[TEnumValueType]TEnumType
}

func (ec *EnumContainer[TEnumType, TEnumValueType]) Has(enum TEnumValueType) bool {
	if _, ok := ec.Map[enum]; ok {
		return true
	}

	return false
}

type Enum[T any] struct {
	Name        string
	Description string
	Value       T
}

func (e Enum[T]) String() string {
	return e.Name
}
