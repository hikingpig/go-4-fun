1. to run all benchmark
go test -bench=.

or only one benchmark:
go test -bench=BenchmarkIsPalindrome

to not run the other Test, only the bench:
go test -bench=BenchmarkIsPalindrome -run "BenchmarkIsPalindrome"

2. using slice,because the underlying array needs to be changed, we must allocate memory.
using array, we only need to allocate memory just once
to see memory allocation use benchmem flag

go test -bench=BenchmarkIsPalindrome -run BenchmarkIsPalindrome -benchmem