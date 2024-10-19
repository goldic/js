# js - Lightweight JSON Utilities for Go

`js` is a lightweight Go package designed to simplify working with JSON. It provides utility types and functions to serialize, deserialize, and manipulate JSON objects, arrays, and values seamlessly.

## Features

- **JSON Parsing and Encoding**: Easily parse and encode JSON data from files, strings, and streams.
- **Flexible Value Handling**: Work with dynamic values of any type, including objects and arrays, without boilerplate code.
- **Utility Functions**: Simplify common operations like comparisons, type conversions, and JSON formatting.
- **Array and Object Manipulation**: Perform common operations on arrays and objects, such as filtering, mapping, and extending.

## Installation

To install the package, use `go get`:

```bash
go get github.com/goldic/js
```

## Usage

### Parsing JSON

You can easily parse JSON data from strings, files, or streams using `Parse` or `ReadValue`:

```go
import (
  "encoding/json"
  "fmt"
  "log"
)

type User struct {
  Name string json:"name"
  Age  int    json:"age"
}

type Response struct {
  Data []User json:"data"
  Ok   bool   json:"ok"
}

func main() {
  response := []byte({"data": [{"name": "Alice", "age": 30}, {"name": "Bob", "age": 25}, {"name": "Charlie", "age": 35}], "ok": true})
  var res Response
  if err := json.Unmarshal(response, &res); err != nil {
  log.Fatal(err)
  }
  
  `", user.Name, user.Age)
  }
  }
  }}`)
  user, err := js.Parse(data)
  if err != nil {
    panic(err)
  }
  
  }
```

### Working with Objects

You can create and manipulate JSON objects easily:

```go
obj := js.NewObject(map[string]any{
    "name": "Alice",
    "age": 30,
})

obj.Set("city", "Wonderland")
fmt.Println(obj.String()) // Output: {"name":"Alice","age":30,"city":"Wonderland"}
```

### Working with Arrays

Easily work with JSON arrays:

```go
arr := js.NewArray("Apple", "Banana", "Cherry")
arr.Push("Date")
fmt.Println(arr.String()) // Output: ["Apple","Banana","Cherry","Date"]
```

### Encoding and Writing JSON

Serialize any value to JSON or write it to a file or stream:

```go
data := js.NewValue(map[string]any{"name": "Bob", "age": 25})
fmt.Println(data.JSON()) // Output: {"name":"Bob","age":25}

err := js.MarshalToFile("data.json", data)
if err != nil {
    panic(err)
}
```

### Utility Functions

- **Comparisons**: Easily compare values with `Cmp`, `Equal`.
- **Type Conversions**: Convert to strings, numbers, and more with `ToStr`, `ToNum`, `ToInt`.
- **Logical Operations**: Simplify conditional logic with `Or` and `And`.

```go
v1 := js.NewValue(10)
v2 := js.NewValue("10")

fmt.Println(v1.Equal(v2)) // Output: true
fmt.Println(js.Cmp(15, 10)) // Output: 1
```

## API Documentation for js Package

### type `Object map[string]any`

Represents a JSON object (`map[string]any`). It is a flexible and dynamic structure to hold key-value pairs. You can set, get, delete, and iterate over the keys.

#### Methods

- **`NewObject(v any) Object`**  
  Creates a new `Object` from any valid Go value (typically a map).

- **`Set(key string, value any)`**  
  Adds or updates a key-value pair in the object.

- **`Get(key string) Value`**  
  Retrieves the value associated with a given key as a `Value`. If the key doesn't exist, it returns an empty `Value`.

- **`Delete(key string)`**  
  Removes a key-value pair from the object.

- **`Has(key string) bool`**  
  Checks whether a key exists in the object.

- **`Keys() []string`**  
  Returns all the keys in the object as a slice of strings.

- **`String() string`**  
  Converts the entire object to a JSON string.

- **`Len() int`**  
  Returns the number of key-value pairs in the object.

- **`IsNull() bool`**  
  Returns true if the object is nil.

- **`Unmarshal(v any) error`**  
  Deserializes the object’s JSON representation into the provided Go value.

### type`Array []any`

Represents a JSON array (`[]any`). It is designed to manage a list of elements, providing various utility methods to manipulate and retrieve data.

#### Methods

- **`NewArray(value ...any) Array`**  
  Creates a new `Array` from the provided elements.

- **`ToArray[T any](vv ...T) Array`**  
  Converts any slice of a specific type to an `Array`.

- **`Push(v any)`**  
  Appends a value to the array.

- **`Eq(i int) Value`**  
  Returns the element at the given index as a `Value`.

- **`ForEach(fn func(Value))`**  
  Iterates over each element in the array and applies the provided function to it.

- **`ForEachObject(fn func(Object))`**  
  Iterates over each element in the array and applies the provided function, treating each element as an `Object`.

- **`FindObject(fn func(Object) bool) Object`**  
  Returns the first object in the array that matches the condition specified by the function.

- **`FindObjectBy(param string, val any) Object`**  
  Returns the first object where the key `param` matches the value `val`.

- **`Filter(fn func(Value) bool) Array`**  
  Filters the array based on a condition and returns a new array with elements that satisfy the condition.

- **`Map(fn func(Value) any) Array`**  
  Applies a function to each element in the array and returns a new array with the results.

- **`Join(sep string) string`**  
  Joins the array elements into a single string, using the provided separator.

- **`Objects() []Object`**  
  Converts all elements of the array to objects and returns them as a slice.

- **`Strings() []string`**  
  Converts all elements of the array to strings.

- **`Ints() []int`**  
  Converts all elements of the array to integers.

- **`Nums() []float64`**  
  Converts all elements of the array to floating-point numbers.

- **`Len() int`**  
  Returns the number of elements in the array.

- **`IsNull() bool`**  
  Returns true if the array is nil.

- **`Bytes() []byte`**  
  Returns the array as a JSON-encoded byte slice.

- **`String() string`**  
  Converts the entire array to a JSON string.

- **`Unmarshal(v any) error`**  
  Deserializes the array’s JSON representation into the provided Go value.

### type `Value struct{}`

Represents a flexible type that can hold any JSON value. It abstracts away the specific JSON type, allowing you to work with any JSON data as a `Value`. It provides methods to convert the value into specific types like strings, numbers, or arrays.

#### Methods

- **`NewValue(v any) Value`**  
  Creates a new `Value` from any Go value.

- **`Parse(data []byte) (Value, error)`**  
  Parses a JSON byte slice and returns a `Value`.

- **`MustParse(data []byte) Value`**  
  Parses a JSON byte slice and returns a `Value`, panicking on error.

- **`ReadValue(r io.Reader) (Value, error)`**  
  Reads JSON from an `io.Reader` and returns a `Value`.

- **`Value() any`**  
  Returns the underlying Go value.

- **`String() string`**  
  Converts the value to a string.

- **`Int() int`**  
  Converts the value to an integer.

- **`Int64() int64`**  
  Converts the value to a 64-bit integer.

- **`Uint64() uint64`**  
  Converts the value to an unsigned 64-bit integer.

- **`Float64() float64`**  
  Converts the value to a floating-point number.

- **`Bool() bool`**  
  Converts the value to a boolean.

- **`Array() Array`**  
  Converts the value to an `Array`.

- **`Object() Object`**  
  Converts the value to an `Object`.

- **`Objects() []Object`**  
  Converts the value to a slice of `Object` if it's an array.

- **`IsNull() bool`**  
  Returns true if the value is nil.

- **`IsNum() bool`**  
  Returns true if the value is a number.

- **`IsObject() bool`**  
  Returns true if the value is an object.

- **`IsArray() bool`**  
  Returns true if the value is an array.

- **`Equal(b any) bool`**  
  Checks whether the current value is equal to another value.

- **`Cmp(b any) int`**  
  Compares the value with another value, returning -1, 0, or 1.

- **`Empty() bool`**  
  Checks if the value is empty (e.g., empty string, nil object, etc.).

- **`Bytes() []byte`**  
  Converts the value to a JSON-encoded byte slice.

- **`JSON() string`**  
  Converts the value to a JSON string.

- **`MarshalJSON() ([]byte, error)`**  
  Marshals the value to a JSON byte slice.

- **`UnmarshalJSON(data []byte) error`**  
  Unmarshals a JSON byte slice into the value.

- **`MarshalTo(val any) error`**  
  Unmarshals the JSON representation of the value into the provided Go value.

### Utility Functions

- **`Marshal(v any) ([]byte, error)`**  
  Marshals any Go value into JSON.

- **`Encode(v any) string`**  
  Encodes any Go value into a JSON string.

- **`IndentEncode(v any) string`**  
  Encodes any Go value into a prettified JSON string.

- **`MarshalToFile(filename string, v any) error`**  
  Marshals a value and writes it to a file as JSON.

- **`MarshalIndentToFile(filename string, v any) error`**  
  Marshals a value as prettified JSON and writes it to a file.

- **`UnmarshalFile(filename string, v any) error`**  
  Reads JSON from a file and unmarshals it into a Go value.

- **`ParseFile(filename string) (Value, error)`**  
  Parses a JSON file into a `Value`.

- **`Write(w io.Writer, v any) error`**  
  Encodes a Go value as JSON and writes it to an `io.Writer`.

- **`Cmp(a, b any) int`**  
  Compares two values, returning -1, 0, or 1.

- **`ToStr(v any) string`**  
  Converts a value to a string.

- **`ToNum(v any) float64`**  
  Converts a value to a floating-point number.

- **`ToInt(v any) int`**  
  Converts a value to an integer.

- **`ToUint64(v any) uint64`**  
  Converts a value to an unsigned 64-bit integer.

- **`IsEmpty(v any) bool`**  
  Checks if a value is empty.

- **`Or(v ...any) any`**  
  Returns the first non-empty value from a list of values.

- **`And(v ...any) bool`**  
  Returns true if all values are non-empty.


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
