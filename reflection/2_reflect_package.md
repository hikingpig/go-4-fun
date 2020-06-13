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