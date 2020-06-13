Reflection works at `run time`, without knowing variable types at compile time.
- update variables and inspect their values at run time
- call their methods
- apply the operations intrinsic to their representation

Reflection is used in:
- string formatting of fmt
- protocol encoding in `encoding/json` and `encoding/xml`
- template mechanism `text/template`, `html/template`

Reflection increases expressiveness of language but complex and are not exposed in the packages's API

Reflection is used to build a function can deal uniformly with types that:
  - NOT satisfy a common interfaces
  - NOT have a known representation
  - NOT exist an the time the function is designed

Naive Example (1_naive_sprintf.go): Sprintf: use switch case to get
  - type switch -> Type has a String method or not
  - add switch to test value's type: string, int, bool, etc

There are infinite number of types. The switch just can not handle all.