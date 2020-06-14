1. methods to inspect in reflect.Values
  - v.Len -> apply only for slice, array, map, Chan and String
  - v.Index : return element i. apply only for slice, array, string
  - v.NumField: apply only to struct
  - v.Type().Field(i): return a reflect.StructField of the ith field. apply only for struct. reflect.StructField.Name is the name of the struct field
  - v.MapKeys: return all keys in a map, unsorted
  - v.MapIndex: return value for a key
  - v.IsNil: apply for func, chan, map, interface, slice, pointer
  - v.Elem: apply for pointer or interface. returns value the interface contains or that the ptr points to.

2. execution flow of x_5
- feed the name "strangelove" and the strangelove movie to Display
- Display feed "strangelove" and `strangelove` to display
- display detects this is a struct
  - it loops over 7 struct fields
    - i = 0, field name is `Title`, new path is "strangelove.Title". call display with new path and Field value
      - the display will detect it as string, call formatAtom to print it in quotes
      - exit the first recursive loop
    - i = 1, field name is `Subtitle`, field value is string. similar to i = 0
    - i = 2. field name is `Year`, field value is int. similarly
    - i = 3. similarly
    - i = 4, field name is `Actor`. field value is map.
      - get all the map keys and loop over each key:
        - the new path is "strangelove.Actor[key]"
        - display is called with new path and element of map
          - it is string -> quoted
        - end the recursive in map element
      - end loop
      - end recursive for i = 4
    
    - i = 5, field name is `Oscars`, field value is slice
      - loop over the len of slice
        - i = 0
          - new path "strangelove.Oscars[0]", value is string.
          - display is called with new path and string
            - quoted
          - end recursive for i = 0
        - i = 1: similarly
      - end loop, end recursive for i = 5
    - i = 6, field name `Sequel`, field value is pointer.
      - the pointer is nil. print it out

3. the path is added up during recursion until the final, innermost display causes it to be printed out:
- strangelove -> strangelove.Oscars -> strangelove.Oscars[0] -> printed as field value is of kind `String`

4. for the pointer as in os.Stderr, we continue to call display with the value the pointer pointed at. and if it is still a pointer, we will continue to get the underlying value. that way, we can see the entire fields nested inside multiple layers of pointers
- this is a very useful technique 

5. even the private field will be displayed with reflection, like 
`(*(*Os.Stderr).file).pfd.IsStream = true`

6. difference of Display("i", i) with
Display("&i", &i)
- the i is an interface, but we call reflect.ValueOf(i) when we passed it to display, it is an int and we jump straight to basic types case
- the &i is an pointer, which nests an interface and that contains an int. display has to go through all 3 levels

7. cycle struct can causes display run forever!
- fmt.Sprint will only print the address of pointer so it will not stuck
- but for map or struct that nested itself, fmt.Sprint will also be stuck

