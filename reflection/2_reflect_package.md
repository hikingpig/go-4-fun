reflect pacage defines two important types: reflect.Type and reflect.Value

- Type:
  - represents Go type
  - has methods to inspect:
    - fields of struct
    - params of functions
  - is implemented only in `type descriptor`, which also identifies dynamic type of interface value
  - reflect.TypeOf: accept an interface and return its dynamic type as a reflect.Type. (see x_2_type_of.go)
    - the returned dynamic type is always a concrete type (the implementation of an interface)
    - %T used in fmt.Fprintf used reflect.TypeOf internally to get the concrete type for debugging

- Value
  - hold a value of any type
  - reflect.ValueOf: (see `x_3_value_of.go`)
    - accepts interface{} and return reflect.Value of the interface's dynamic value
    - result is always concrete
    - is used internally of %v in fmt.Printf
    - the String() method only return the type of value unless the value is string

     * so WHY Println and Printf doesn't use this String() method?
     * fmt package `TREAT Values specially`. It doesn't use the String() method but print the underlying concrete value!

    - The Type() method return reflect.Type of the value
  - reflect.Value.Interface is the reverse of reflect.ValueOf. it returns an interface{} that holds THE SAME concrete value as reflect.Value's
    - empty interface hides the underlying type while reflect.Value has many methods to inspect it

  - Kind() method: see x_4_kind.go
    - there are only finite number of kinds, which contains all types:
      - basics: Bool, String, numbers
      - aggregates: Array, Struct
      - references: Chan, Struct, Ptr, Slice, Map
      - Interface
      - Invalid
    - there are methods to inspect the type of underlying value
      - v.Int() will panic if value is not of int types
      - v.Uint() similarly
      - v.Bool()
      - v.String: for value not of type string, it prints <Type Value> like <[]int64 Value>
    - we use v.Type().String() to get the real type
