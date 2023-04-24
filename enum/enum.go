package enum

type EnumContainer[TEnumType any, TEnumValueType comparable] struct {
	Values      []TEnumType
	Map         map[TEnumValueType]TEnumType
	compareFunc func(a, b any) bool
}

func NewEnumContainer[TEnumType any, TEnumValueType comparable](
	valuesMap map[TEnumValueType]TEnumType,
	compareFunc ...func(a, b any) bool,
) *EnumContainer[TEnumType, TEnumValueType] {

	cont := &EnumContainer[TEnumType, TEnumValueType]{
		Values:      make([]TEnumType, 0),
		Map:         valuesMap,
		compareFunc: nil,
	}

	if len(compareFunc) > 0 {
		cont.compareFunc = compareFunc[0]
	}

	for _, v := range valuesMap {
		cont.Values = append(cont.Values, v)
	}

	return cont
}

func (ec *EnumContainer[TEnumType, TEnumValueType]) Has(enum TEnumValueType) bool {
	if ec.compareFunc != nil {
		for _, v := range ec.Values {
			if ec.compareFunc(v, enum) {
				return true
			}
		}
		return false
	}

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
