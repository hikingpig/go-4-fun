1. data is stored and loaded more efficiently if their addresses follow some patterns:
  - 2-byte values (int16), addresses are even numbers (2 * x)
  - 4-byte values (int32), addresses are multiples of 4 (4 * x)
  - 8-byte values (float64, uint64), addresses are mulitples of 8 (8 * x)
  - higher values doesn't required patterns

- these patterns are called alignement

2. aggregate types such as struct or array, the size may be larger than the sum of its fields. there are "holes" added to make the alignment for storage of struct of array. Holes are unused addresses in memory

