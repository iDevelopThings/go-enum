# go-enum

go-enum is a simple and powerful Golang code generator for creating type-safe and human-readable enums. It allows you to define your enums using either struct tags or a simple comment structure. The
generated code provides an easy-to-use API for interacting with your enums.

## Installation

To install the idevelopthings/go-enum package, simply run:

```bash
go get -u github.com/iDevelopThings/go-enum
```

## Usage

Create a Go file containing your enum definitions using either of the two supported formats: struct tags or comments. Then use the go:generate directive to run the go-enum code generator.

### Struct Tags Format

```go
//go:generate go run github.com/iDevelopThings/go-enum -types=MyEnumName,MyOtherEnum -output=my_enums.gen.go -input=my_enums.go

type MyEnumName struct {
    SomeEnumValue string `name:"Special" description:"It's a special value with a special name"`
    // It's another special value... with a description defined via comment
    AnotherSpecialValue string
}
```

### Comments Format

```go
// Start:MyOtherEnum

// Name: SomeEnumValue
// It's a special value with a special name

// Name: AnotherSpecialValue
// It's another special value... with a description defined via comment

// End:MyOtherEnum
```

# Generated API

The generated file contains enums with a consistent API, making it simple for you to use the generated enums in your code.

## Enum Type
Each generated enum type (e.g., MyName and MyOther) has an embedded Enum[string] field containing the enum value, name, and description. The enum type also has a GetAll() method that returns a
container object containing all enum values.

## Enum Container
The container object (e.g., MyNamesContainer and MyOthersContainer) is a struct that holds a map and a slice of all enum values. It also contains named fields for each enum value for easy access.

## Accessing Enum Values
You can access the enum values through the container object, for example:

```go
myEnumValue := MyNames.SomeEnumValue
myOtherEnumValue := MyOthers.AnotherSpecialValue
```

## Enum Properties

You can access the properties of each enum value, such as its name, value, and description:

```go
fmt.Println(myEnumValue.Name) // Output: "SomeEnumValue"
fmt.Println(myEnumValue.Value) // Output: "Special"
fmt.Println(myEnumValue.Description) // Output: "It's a special value with a special name"
```

## Enum Container Properties

You can access the slice and map of all enum values using the EnumContainer:

```go
fmt.Println(MyNames.EnumContainer.Values) // Output: A slice of all MyName enum values
fmt.Println(MyNames.EnumContainer.Map) // Output: A map of all MyName enum values, keyed by their names
```

## Iterate Over Enum Values

You can iterate over all enum values using a loop:

```go
for _, myEnum := range MyNames.EnumContainer.Values {
    fmt.Println(myEnum.Name, myEnum.Value, myEnum.Description)
}
```
