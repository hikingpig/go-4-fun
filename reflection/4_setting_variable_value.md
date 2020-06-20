1. variable = `addressable storage location` that contain a value and its value may be updated through an address

- for refect.Value
```go
reflect.ValueOf(2) // not addressable, only a copy of 2
reflect.ValueOf(x) // only a copy of x's value
reflect.ValueOf(&x) // only a copy of x's address

reflect.ValueOf(&x).Elem() // ADDRESSABLE. it returns the variable of x's address!
```
- checking addressablility of a `reflect.Value` with:
```go
a.CanAddr()
```

- access a value through pointer -> always addressable
```go
// e is an array
e[i] // implicitly use a pointer to access element -> addressable
reflect.ValueOf(e).Index(i) // also addressable
```
- 3 steps to recover a variable from addressable reflect.Value
- Set method
```go
x := 2
d := reflect.ValueOf(&x).Elem()
/*
1. call Addr() to get the pointer to that variable
2. call Interface() to get an interface{} containing value of pointer
3. assert type of the variable to retrieve contents of interface as ordinary pointer
*/
px := d.Addr().Interface().(*int)
*px = 3
fmt.Println(x)

/*
update the variable referred by addressable reflect.Value directly
*/
d.Set(reflect.ValueOf(4))
fmt.Println(x)

// make sure assignability of variable
d.Set(reflect.ValueOf(int64(5))) // this will panic

// calling Set on non-addressable reflect.Value will also panic
e := reflect.ValueOf(x)
e.Set(reflect.ValueOf(3)) // PANIC

// SetInt, SetUint, SetFloat: for specific type
d.SetInt(3)

// SetInt is applied for typed variable, not interface. Interface{} can accept any type!
var y interface{}
ry := reflect.ValueOf(&y).Elem()
ry.SetInt(2) // PANIC
ry.Set(reflect.ValueOf(3)) // OK

ry.SetString("hello") // PANIC
ry.Set(reflect.ValueOf("hello")) // OK

// reflect can read private fields but can not update them
stdout := reflect.ValueOf(os.Stdout).Elem()

fd := stdout.FieldByName("fd")
fmt.Println(fd.Int())

fd.SetInt(2) // PANIC: unexported fields

// CanSet method will check if we can set an addressable reflect.Value
fmt.Println(fd.CanAddr(), fd.CanSet()) // true, false
```

