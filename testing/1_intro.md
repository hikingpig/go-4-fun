1. files with names _test.go are not part of the package and ignored by go build

2. 3 kinds of functions are treated specially:
  - tests
    - name begin with Test
  - benchmarks
    - name begin with Benchmark
    - measure performance
  - examples
    - name begin with Example
    - machine-checked documentation

3. go test will:
  - scan the _test.go files
  - generate temporary main package to call them all
  - builds and runs it
  - report results
  - clean up

4. Test func
```go
func TestSin(t *testing.T) {}
```
the t param provides methods for:
  - reporting test failures
  - logging additional info

5. must remove draft.go, it can not be empty, nor having any illegal syntax. the go test parse this file, all file in the folder!

6. go test -v print the name and execution time of the tests we run
- so we can see which tests are slow and which are quick

7. we can choose which test to run using -run
go test -v -run "French|Canal"
- run takes a regular expression. We can use French without the quotes. But for expression with | we use quotes to prevent misunderstanding as piping

